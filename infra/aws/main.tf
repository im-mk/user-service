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
    key    = "code-service-ci.tfstate"
    region = "eu-west-2"
  }
}

provider "aws" {
  region = var.aws_region
}

resource "aws_ecr_repository" "go_app_ecr" {
  name = "${var.app_name}-repo"
}
