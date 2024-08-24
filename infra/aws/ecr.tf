resource "aws_ecr_repository" "ecr" {
  name = "${var.app_name}-ecr"
  tags = {
    Application = var.app_name
    Project     = var.project
    Environment = var.app_environment
  }
}
