variable "aws_region" {
  default = "eu-west-2"
  type    = string
}

variable "app_name" {
  default = "user-service"
  type    = string
}

variable "repo_name" {
  default = "user-service"
  type    = string
}


variable "artifacts_bucket" {
  default = "im-mk-app-artifacts"
  type    = string
}

variable "code_branch" {
  default = "feature/boilerplate"
  type    = string
}
