output "repository_url" {
  description = "The URL of the created ECR repository."
  value       = aws_ecr_repository.ecr.repository_url
}

output "alb_hostname" {
  value = aws_alb.main.dns_name
}
