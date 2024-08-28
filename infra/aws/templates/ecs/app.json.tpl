[
  {
    "name": "${app_name}",
    "image": "${app_image}",
    "cpu": ${fargate_cpu},
    "memory": ${fargate_memory},
    "networkMode": "awsvpc",
    "logConfiguration": {
        "logDriver": "awslogs",
        "options": {
          "awslogs-group": "/ecs/${app_name}",
          "awslogs-region": "${aws_region}",
          "awslogs-stream-prefix": "ecs"
        }
    },
    "portMappings": [
      {
        "containerPort": ${app_port},
        "hostPort": ${app_port}
      }
    ],
    "healthCheck": {
      "command": [
        "CMD-SHELL",
        "curl -f http://localhost:${app_port}/${health_path} || exit 1"
      ],
      "interval": 30,
      "timeout": 5,
      "retries": 3,
      "startPeriod": 0
    },
    "secrets": [
      {
        "name": "DB_PASSWORD",
        "valueFrom": "arn:aws:ssm:${aws_region}:${account_id}:parameter/user-service-db/db-password"
      },
      {
        "name": "DB_USER",
        "valueFrom": "${db_username}"
      },
      {
        "name": "DB_HOST",
        "valueFrom": "${db_host}"
      },
      {
        "name": "DB_DBNAME",
        "valueFrom": "${db_name}"
      }
    ]
  }
]