variable "aws_region" {
  default = "eu-west-2"
  type    = string
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
