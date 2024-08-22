output "repository_url" {
  description = "The URL of the created ECR repository."
  value       = aws_ecr_repository.go_app_ecr.repository_url
}
