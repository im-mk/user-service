data "aws_ssm_parameter" "db_password" {
  name = "/user-service-db/db-password"
}

resource "aws_db_subnet_group" "user_service_db_subnet_group" {
  name       = "${var.app_name}-db-subnet-group"
  subnet_ids = local.private_subnet_ids

  tags = {
    Name        = "${var.app_name}-db-subnet-group"
    Application = var.app_name
    Project     = var.project
    Environment = var.app_environment
  }
}

resource "aws_db_instance" "user_service_db" {
  allocated_storage    = var.db_allocated_storage
  engine               = "postgres"
  engine_version       = "16.1"
  instance_class       = var.db_instance_class
  db_name              = var.db_name
  identifier           = var.db_identifier
  username             = var.db_username
  password             = data.aws_ssm_parameter.db_password.value
  db_subnet_group_name = aws_db_subnet_group.user_service_db_subnet_group.name
  vpc_security_group_ids = [
    aws_security_group.user_service_db.id
  ]
  skip_final_snapshot = true

  tags = {
    Name        = "${var.app_name}-db"
    Application = var.app_name
    Project     = var.project
    Environment = var.app_environment
  }
}

resource "aws_security_group" "user_service_db" {
  name        = "${var.app_name}-db-sg"
  description = "Security group for PostgreSQL RDS"
  vpc_id      = local.vpc_id

  ingress {
    from_port   = 5432
    to_port     = 5432
    protocol    = "tcp"
    cidr_blocks = ["10.0.0.0/16"]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = {
    Name        = "${var.app_name}-db-sg"
    Application = var.app_name
    Project     = var.project
    Environment = var.app_environment
  }
}
