terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 4.0"
    }
  }
  required_version = ">= 0.14"
  backend "s3" {
    bucket = "internal-tfstate"
    key    = "tf/code-service.tfstate"
    region = "eu-west-2"
  }
}

provider "aws" {
  region = var.aws_region
}

data "aws_secretsmanager_secret_version" "github_token" {
  secret_id = "github/token"
}

resource "aws_ecr_repository" "go_app_ecr" {
  name = "${var.app_name}-repo"
}

resource "aws_iam_role" "codebuild_role" {
  name = "${var.app_name}-codebuild-role"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Action = "sts:AssumeRole"
        Effect = "Allow"
        Principal = {
          Service = "codebuild.amazonaws.com"
        }
      }
    ]
  })

  managed_policy_arns = [
    "arn:aws:iam::aws:policy/AmazonEC2ContainerRegistryPowerUser",
    "arn:aws:iam::aws:policy/AWSCodeBuildDeveloperAccess",
  ]
}

resource "aws_iam_role_policy" "codebuild_cloudwatch_logs_access" {
  name = "codebuild-cloudwatch-logs-access"
  role = aws_iam_role.codebuild_role.name

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect = "Allow",
        Action = [
          "logs:CreateLogGroup",
          "logs:CreateLogStream",
          "logs:PutLogEvents"
        ],
        Resource = "*"
      }
    ]
  })
}

resource "aws_iam_role_policy" "codebuild_s3_access" {
  name = "codepipeline-s3-access"
  role = aws_iam_role.codebuild_role.name

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect = "Allow"
        Action = [
          "s3:GetObject",
          "s3:GetObjectVersion",
          "s3:ListBucket",
        ]
        Resource = [
          "${aws_s3_bucket.codepipeline_bucket.arn}",
          "${aws_s3_bucket.codepipeline_bucket.arn}/*"
        ]
      }
    ]
  })
}

resource "aws_codebuild_project" "go_app_build" {
  name          = "${var.app_name}-build"
  description   = "Build project for Go application"
  service_role  = aws_iam_role.codebuild_role.arn
  build_timeout = 5

  artifacts {
    type = "NO_ARTIFACTS"
  }

  source {
    type      = "GITHUB"
    location  = "https://github.com/im-mk/user-service.git"
    buildspec = <<EOF
version: 0.2

phases:
  install:
    commands:
      - echo Installing Go dependencies...
      - cd src 
      - go mod download
  pre_build:
    commands:
      - echo Logging in to Amazon ECR...
      - aws ecr get-login-password --region ${var.aws_region} | docker login --username AWS --password-stdin ${aws_ecr_repository.go_app_ecr.repository_url}
  build:
    commands:
      - echo Building the Docker image...
      - docker build -t ${aws_ecr_repository.go_app_ecr.repository_url}:latest .
      - docker tag ${aws_ecr_repository.go_app_ecr.repository_url}:latest ${aws_ecr_repository.go_app_ecr.repository_url}:$CODEBUILD_RESOLVED_SOURCE_VERSION
  post_build:
    commands:
      - echo Pushing the Docker image...
      - docker push ${aws_ecr_repository.go_app_ecr.repository_url}:latest
      - docker push ${aws_ecr_repository.go_app_ecr.repository_url}:$CODEBUILD_RESOLVED_SOURCE_VERSION
EOF
  }

  environment {
    compute_type    = "BUILD_GENERAL1_SMALL"
    image           = "aws/codebuild/standard:5.0"
    type            = "LINUX_CONTAINER"
    privileged_mode = true
    environment_variable {
      name  = "AWS_REGION"
      value = var.aws_region
    }
  }

  logs_config {
    cloudwatch_logs {
      group_name  = "go-app-build-logs"
      stream_name = "go-app-build"
    }
  }
}

resource "aws_iam_role" "codepipeline_role" {
  name = "${var.app_name}-codepipeline-role"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Action = "sts:AssumeRole"
        Effect = "Allow"
        Principal = {
          Service = "codepipeline.amazonaws.com"
        }
      }
    ]
  })

  managed_policy_arns = [
    "arn:aws:iam::aws:policy/AWSCodePipeline_FullAccess",
    "arn:aws:iam::aws:policy/AWSCodeBuildAdminAccess",
  ]
}

resource "aws_iam_role_policy" "codepipeline_s3_access" {
  name = "codepipeline-s3-access"
  role = aws_iam_role.codepipeline_role.name

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect = "Allow"
        Action = [
          "s3:GetObject",
          "s3:GetObjectVersion",
          "s3:ListBucket",
          "s3:PutObject"
        ]
        Resource = [
          "${aws_s3_bucket.codepipeline_bucket.arn}",
          "${aws_s3_bucket.codepipeline_bucket.arn}/*"
        ]
      }
    ]
  })
}

resource "aws_codepipeline" "go_app_pipeline" {
  name     = "${var.app_name}-pipeline"
  role_arn = aws_iam_role.codepipeline_role.arn

  artifact_store {
    location = aws_s3_bucket.codepipeline_bucket.bucket
    type     = "S3"
  }

  stage {
    name = "Source"

    action {
      name             = "Source"
      category         = "Source"
      owner            = "ThirdParty"
      provider         = "GitHub"
      version          = "1"
      output_artifacts = ["source_output"]

      configuration = {
        Owner      = "im-mk"
        Repo       = var.repo_name
        Branch     = var.code_branch
        OAuthToken = data.aws_secretsmanager_secret_version.github_token.secret_string
      }
    }
  }

  stage {
    name = "Build"

    action {
      name             = "Build"
      category         = "Build"
      owner            = "AWS"
      provider         = "CodeBuild"
      input_artifacts  = ["source_output"]
      output_artifacts = ["build_output"]
      version          = "1"

      configuration = {
        ProjectName = aws_codebuild_project.go_app_build.name
      }
    }
  }
}

resource "aws_s3_bucket" "codepipeline_bucket" {
  bucket = var.artifacts_bucket
}

output "ecr_repository_url" {
  value = aws_ecr_repository.go_app_ecr.repository_url
}
