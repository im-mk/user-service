resource "aws_ecs_cluster" "main" {
  name = "${var.app_name}-${var.app_environment}-cluster"
  tags = {
    Name        = "${var.app_name}-ecs"
    Application = var.app_name
    Project     = var.project
    Environment = var.app_environment
  }
}

data "template_file" "app" {
  template = file("./templates/ecs/app.json.tpl")

  vars = {
    app_image      = "${aws_ecr_repository.ecr.repository_url}:${var.app_image}"
    app_port       = var.container_port
    fargate_cpu    = var.fargate_cpu
    fargate_memory = var.fargate_memory
    aws_region     = var.aws_region
    app_name       = var.app_name
  }
}

resource "aws_ecs_task_definition" "app" {
  family                   = "${var.app_name}-task"
  execution_role_arn       = aws_iam_role.ecs_task_execution_role.arn
  network_mode             = "awsvpc"
  requires_compatibilities = ["FARGATE"]
  cpu                      = var.fargate_cpu
  memory                   = var.fargate_memory
  container_definitions    = data.template_file.app.rendered
  tags = {
    Name        = "${var.app_name}-ecs-td"
    Application = var.app_name
    Project     = var.project
    Environment = var.app_environment
  }
}

resource "aws_ecs_service" "main" {
  name            = "${var.app_name}-service"
  cluster         = aws_ecs_cluster.main.id
  task_definition = aws_ecs_task_definition.app.arn
  desired_count   = var.az_count
  launch_type     = "FARGATE"

  network_configuration {
    security_groups  = [aws_security_group.ecs_tasks.id]
    subnets          = local.private_subnet_ids
    assign_public_ip = true
  }

  load_balancer {
    target_group_arn = aws_alb_target_group.app.id
    container_name   = var.app_name
    container_port   = var.container_port
  }
}
