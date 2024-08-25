# Set up CloudWatch group and log stream and retain logs for 30 days
resource "aws_cloudwatch_log_group" "cb_log_group" {
  name              = "/ecs/${var.app_name}"
  retention_in_days = 30

  tags = {
    Name        = "${var.app_name}-log-group"
    Application = var.app_name
    Project     = var.project
    Environment = var.app_environment
  }
}

resource "aws_cloudwatch_log_stream" "cb_log_stream" {
  name           = "${var.app_name}-log-stream"
  log_group_name = aws_cloudwatch_log_group.cb_log_group.name
}
