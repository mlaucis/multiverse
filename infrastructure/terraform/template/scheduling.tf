resource "aws_ecs_cluster" "service" {
  name = "service"
}

resource "aws_ecs_service" "gateway-http" {
  cluster         = "${aws_ecs_cluster.service.id}"
  depends_on      = [
    "aws_iam_instance_profile.ecs-agent-profile",
    "aws_db_instance.master",
    "aws_elasticache_cluster.ratelimiter",
  ]
  deployment_maximum_percent          = 200
  deployment_minimum_healthy_percent  = 50
  desired_count   = 2
  iam_role        = "${aws_iam_role.ecs-scheduler.arn}"
  name            = "gateway-http"
  task_definition = "${aws_ecs_task_definition.gateway-http.arn}"

  load_balancer {
    container_name  = "gateway-http"
    container_port  = 8083
    elb_name        = "${aws_elb.gateway-http.id}"
  }
}

resource "aws_ecs_task_definition" "gateway-http" {
  family                = "gateway-http"
  container_definitions = <<EOF
[
  {
    "command": [
      "./gateway-http"
    ],
    "cpu": 512,
    "dnsSearchDomains": [
      "${var.env}.${var.region}"
    ],
    "essential": true,
    "image": "775034650473.dkr.ecr.us-east-1.amazonaws.com/gateway-http:${var.version.gateway-http}",
    "logConfiguration": {
      "logDriver": "syslog"
    },
    "memory": 2048,
    "name": "gateway-http",
    "portMappings": [
      {
        "containerPort": 8083,
        "hostPort": 8083
      },
      {
        "containerPort": 9000,
        "hostPort": 9000
      }
    ],
    "readonlyRootFilesystem": true,
    "workingDirectory": "/tapglue/"
  }
]
EOF
}