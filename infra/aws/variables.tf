variable "aws_region" {
  default     = "eu-west-2"
  description = "AWS Region to deploy the application"
  type        = string
}

variable "app_name" {
  default = "user-service"
  type    = string
}

variable "project" {
  default = "clover"
  type    = string
}

variable "app_environment" {
  default = "production"
  type    = string
}

variable "cidr" {
  description = "The CIDR block for the VPC."
  type        = string
  default     = "10.0.0.0/16"
}

variable "az_count" {
  default     = 2
  description = "Number of availability zones"
}

variable "app_port" {
  description = "alb port"
  default     = 80
}

variable "container_port" {
  description = "Port exposed by the docker image to redirect traffic to"
  default     = 8080
}

variable "health_check_path" {
  default = "/"
}

variable "fargate_cpu" {
  description = "Fargate instance CPU units to provision (1 vCPU = 1024 CPU units)"
  default     = "512"
}

variable "fargate_memory" {
  description = "Fargate instance memory to provision (in MiB)"
  default     = "1024"
}

variable "app_image" {
  description = "Docker image to run in the ECS cluster"
  default     = "latest"
}