resource "aws_ecr_repository" "repo" {
  name = "${var.app_name}-ecr"

  tags = {
    Name        = "${var.app_name}-ecr"
    Application = var.app_name
    Project     = var.project
    Environment = var.app_environment
  }
}
