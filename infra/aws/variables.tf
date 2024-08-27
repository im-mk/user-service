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

variable "container_count" {
  default     = 1
  description = "Number of containers"
}

variable "app_port" {
  description = "alb port"
  default     = 80
}

variable "container_port" {
  description = "Port exposed by the docker image to redirect traffic to"
  default     = 8080
}

variable "health_path" {
  default = "health"
}

variable "fargate_cpu" {
  description = "Fargate instance CPU units to provision (1 vCPU = 1024 CPU units)"
  default     = "256"
}

variable "fargate_memory" {
  description = "Fargate instance memory to provision (in MiB)"
  default     = "512"
}

variable "app_image" {
  description = "Docker image to run in the ECS cluster"
  default     = "latest"
}

variable "db_allocated_storage" {
  description = "The amount of storage in gigabytes for the database"
  default     = 5
}

variable "db_instance_class" {
  description = "The instance type of the RDS instance"
  default     = "db.t3.micro"
}

variable "db_name" {
  description = "The name for the PostgreSQL database"
  default     = "userservicedb"
}

variable "db_identifier" {
  description = "The instance id of the PostgreSQL database"
  default     = "user-service-db"
}

variable "db_username" {
  description = "The username for the PostgreSQL database"
  default     = "dbadmin"
}
