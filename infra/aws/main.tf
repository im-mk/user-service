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
    key    = "user-service.tfstate"
    region = "eu-west-2"
  }
}

provider "aws" {
  region = var.aws_region
}

// Get infra details
data "terraform_remote_state" "infra" {
  backend = "s3"

  config = {
    bucket = "internal-tfstate"
    key    = "clover-infra.tfstate"
    region = "eu-west-2"
  }
}

# Define locals to store the outputs
locals {
  private_subnet_ids = data.terraform_remote_state.infra.outputs.private_subnet_ids
  public_subnet_ids  = data.terraform_remote_state.infra.outputs.public_subnet_ids
  vpc_id             = data.terraform_remote_state.infra.outputs.vpc_id
}
