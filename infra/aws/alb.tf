resource "aws_alb" "main" {
  name            = "${var.app_name}-alb"
  subnets         = aws_subnet.public.*.id
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
  vpc_id      = aws_vpc.aws-vpc.id
  target_type = "ip"

  health_check {
    healthy_threshold   = "3"
    interval            = "30"
    protocol            = "HTTP"
    matcher             = "200"
    timeout             = "3"
    path                = var.health_check_path
    unhealthy_threshold = "2"
  }

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
