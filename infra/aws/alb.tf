resource "aws_alb" "main" {
  name            = "${var.app_name}-alb"
  subnets         = local.public_subnet_ids
  security_groups = [aws_security_group.lb.id]
  tags = {
    Name        = "${var.app_name}-alb"
    Application = var.app_name
    Project     = var.project
    Environment = var.app_environment
  }
}

resource "aws_alb_target_group" "app" {
  name        = "${var.app_name}-target-group"
  port        = 80
  protocol    = "HTTP"
  vpc_id      = local.vpc_id
  target_type = "ip"

  tags = {
    Name        = "${var.app_name}-target-group"
    Application = var.app_name
    Project     = var.project
    Environment = var.app_environment
  }
}

# Redirect all traffic from the ALB to the target group
resource "aws_alb_listener" "front_end" {
  load_balancer_arn = aws_alb.main.id
  port              = var.app_port
  protocol          = "HTTP"

  default_action {
    target_group_arn = aws_alb_target_group.app.id
    type             = "forward"
  }

  tags = {
    Name        = "${var.app_name}-alb-listener"
    Application = var.app_name
    Project     = var.project
    Environment = var.app_environment
  }
}
