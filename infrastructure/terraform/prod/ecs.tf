resource "aws_ecr_repository" "dashboard" {
  provider = "aws.us-east-1"
  name     = "dashboard"
}

resource "aws_ecr_repository_policy" "dashboard-deployment" {
  provider = "aws.us-east-1"
  repository = "${aws_ecr_repository.dashboard.name}"
  policy     = <<EOF
{
    "Version": "2008-10-17",
    "Statement": [
        {
            "Sid": "deployment",
            "Effect": "Allow",
            "Principal": {
                "AWS": [
                    "arn:aws:iam::775034650473:root",
                    "arn:aws:iam::775034650473:role/ecsInstance",
                    "arn:aws:iam::775034650473:user/deployer"
                ]
            },
            "Action": [
                "ecr:GetDownloadUrlForLayer",
                "ecr:BatchGetImage",
                "ecr:BatchCheckLayerAvailability"
            ]
        }
    ]
}
EOF
}

resource "aws_ecr_repository" "gateway-http" {
  provider = "aws.us-east-1"
  name     = "gateway-http"
}

resource "aws_ecr_repository_policy" "gateway-http-deployment" {
  provider = "aws.us-east-1"
  repository = "${aws_ecr_repository.gateway-http.name}"
  policy     = <<EOF
{
    "Version": "2008-10-17",
    "Statement": [
        {
            "Sid": "deployment",
            "Effect": "Allow",
            "Principal": {
                "AWS": [
                    "arn:aws:iam::775034650473:root",
                    "arn:aws:iam::775034650473:role/ecsInstance",
                    "arn:aws:iam::775034650473:user/deployer"
                ]
            },
            "Action": [
                "ecr:GetDownloadUrlForLayer",
                "ecr:BatchGetImage",
                "ecr:BatchCheckLayerAvailability"
            ]
        }
    ]
}
EOF
}

resource "aws_ecs_cluster" "service" {
  name = "service"
}

resource "aws_iam_role" "ecsInstance" {
  name               = "ecsInstance"
  path               = "/"
  assume_role_policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": "sts:AssumeRole",
      "Principal": {
        "Service": "ec2.amazonaws.com"
      },
      "Effect": "Allow",
      "Sid": ""
    }
  ]
}
EOF
}

resource "aws_iam_role_policy" "ecsOperations" {
  name   = "ecsOperations"
  role   = "${aws_iam_role.ecsInstance.id}"
  policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Action": [
        "ecs:CreateCluster",
        "ecs:DeregisterContainerInstance",
        "ecs:DiscoverPollEndpoint",
        "ecs:Poll",
        "ecs:RegisterContainerInstance",
        "ecs:StartTelemetrySession",
        "ecs:Submit*",
        "ecr:GetAuthorizationToken",
        "ecr:BatchCheckLayerAvailability",
        "ecr:GetDownloadUrlForLayer",
        "ecr:BatchGetImage"
      ],
      "Resource": [
        "*"
      ]
    }
  ]
}
EOF
}

resource "aws_iam_role" "ecsELB" {
  name               = "ecsELB"
  path               = "/"
  assume_role_policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": "sts:AssumeRole",
      "Principal": {
        "Service": "ecs.amazonaws.com"
      },
      "Effect": "Allow",
      "Sid": ""
    }
  ]
}
EOF
}

resource "aws_iam_role_policy" "ecsELB" {
  name   = "tapglueEcsElbRole"
  role   = "${aws_iam_role.ecsELB.id}"
  policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Action": [
        "elasticloadbalancing:Describe*",
        "elasticloadbalancing:DeregisterInstancesFromLoadBalancer",
        "elasticloadbalancing:RegisterInstancesWithLoadBalancer",
        "ec2:Describe*",
        "ec2:AuthorizeSecurityGroupIngress"
      ],
      "Resource": [
        "*"
      ]
    }
  ]
}
EOF
}

resource "aws_elb" "gateway-http" {
  connection_draining         = true
  connection_draining_timeout = 10
  cross_zone_load_balancing   = true
  idle_timeout                = 30
  name                        = "gateway-http"
  subnets                     = [
    "${aws_subnet.public-a.id}",
    "${aws_subnet.public-b.id}",
  ]
  security_groups             = [
    "${aws_security_group.loadbalancer.id}",
  ]

  access_logs {
    bucket   = "tapglue-logs"
    interval = 5
  }

  listener {
    instance_port      = 8083
    instance_protocol  = "https"

    lb_port            = 443
    lb_protocol        = "https"

    ssl_certificate_id = "${aws_iam_server_certificate.self-signed.arn}"
  }

  health_check {
    healthy_threshold   = 2
    interval            = 5
    timeout             = 2
    target              = "HTTPS:8083/health-45016490610398192"
    unhealthy_threshold = 2
  }

  tags {
    Name = "gateway-http"
  }
}

resource "aws_elb" "dashboard" {
  connection_draining         = true
  connection_draining_timeout = 10
  cross_zone_load_balancing   = true
  idle_timeout                = 30
  name                        = "dashboard"
  subnets                     = [
    "${aws_subnet.public-a.id}",
    "${aws_subnet.public-b.id}",
  ]
  security_groups             = [
    "${aws_security_group.loadbalancer.id}",
  ]

  access_logs {
    bucket   = "tapglue-logs"
    interval = 5
  }

  listener {
    instance_port      = 8081
    instance_protocol  = "https"

    lb_port            = 443
    lb_protocol        = "https"

    ssl_certificate_id = "${aws_iam_server_certificate.self-signed.arn}"
  }

  health_check {
    healthy_threshold   = 2
    interval            = 5
    timeout             = 2
    target              = "HTTPS:8081/"
    unhealthy_threshold = 2
  }

  tags {
    Name = "dashboard"
  }
}

resource "aws_ecs_task_definition" "dashboard" {
  family                = "dashboard"
  container_definitions = <<EOF
[
  {
    "name": "dashboard",
    "image": "775034650473.dkr.ecr.us-east-1.amazonaws.com/dashboard:1717",
    "cpu": 256,
    "memory": 256,
    "essential": true,
    "workingDirectory": "/home/tapglue/releases/",
    "readonlyRootFilesystem": false,
    "privileged": false,
    "portMappings": [
      {
        "containerPort": 443,
        "hostPort": 8081
      }
    ],
    "logConfiguration": {
      "logDriver": "syslog"
    }
  }
]
EOF
}

resource "aws_ecs_service" "dashboard" {
  name            = "dashboard"
  cluster         = "${aws_ecs_cluster.service.id}"
  task_definition = "${aws_ecs_task_definition.dashboard.arn}"
  desired_count   = 2
  iam_role = "${aws_iam_role.ecsELB.arn}"
  depends_on = [
    "aws_iam_role_policy.ecsELB",
  ]

  load_balancer = {
    elb_name = "${aws_elb.dashboard.id}"
    container_name = "dashboard"
    container_port = 443
  }
}

resource "aws_ecs_task_definition" "gateway-http" {
  family                = "gateway-http"
  container_definitions = <<EOF
[
  {
    "name": "gateway-http",
    "image": "775034650473.dkr.ecr.us-east-1.amazonaws.com/gateway-http:1717",
    "cpu": 512,
    "memory": 2048,
    "essential": true,
    "workingDirectory": "/tapglue/",
    "readonlyRootFilesystem": true,
    "privileged": false,
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
    "logConfiguration": {
      "logDriver": "syslog"
    }
  }
]
EOF
}

resource "aws_ecs_service" "gateway-http" {
  name            = "gateway-http"
  cluster         = "${aws_ecs_cluster.service.id}"
  task_definition = "${aws_ecs_task_definition.gateway-http.arn}"
  desired_count   = 2
  iam_role = "${aws_iam_role.ecsELB.arn}"
  depends_on = [
    "aws_iam_role_policy.ecsELB",
  ]

  load_balancer = {
    elb_name = "${aws_elb.gateway-http.id}"
    container_name = "gateway-http"
    container_port = 8083
  }
}

resource "aws_iam_instance_profile" "service" {
  name  = "service"
  roles = [
    "${aws_iam_role.ecsInstance.name}"
  ]
}

resource "aws_launch_configuration" "service" {
  image_id                    = "${var.ami_container}"
  instance_type               = "m4.large"
  associate_public_ip_address = false
  enable_monitoring           = true
  ebs_optimized               = false
  iam_instance_profile        = "${aws_iam_instance_profile.service.name}"

  user_data                   = <<EOF
#!/bin/bash
echo ECS_CLUSTER=service >> /etc/ecs/ecs.config

sudo echo '$WorkDirectory /var/spool/rsyslog # where to place spool files' > /etc/rsyslog.d/22-loggly.conf
sudo echo '$ActionQueueFileName fwdRule1     # unique name prefix for spool files' >> /etc/rsyslog.d/22-loggly.conf
sudo echo '$ActionQueueMaxDiskSpace 100m     # 1gb space limit (use as much as possible)' >> /etc/rsyslog.d/22-loggly.conf
sudo echo '$ActionQueueSaveOnShutdown on     # save messages to disk on shutdown' >> /etc/rsyslog.d/22-loggly.conf
sudo echo '$ActionQueueType LinkedList       # run asynchronously' >> /etc/rsyslog.d/22-loggly.conf
sudo echo '$ActionResumeRetryCount -1        # infinite retries if host is down' >> /etc/rsyslog.d/22-loggly.conf
sudo echo '' >> /etc/rsyslog.d/22-loggly.conf
sudo echo '$template LogglyFormat,"<%pri%>%protocol-version% %timestamp:::date-rfc3339% %HOSTNAME% %app-name% %procid% %msgid% [d2e7097f-25aa-497a-a9e3-d691bd4ec7ab@41058 tag=\"service\"] %msg%\n"' >> /etc/rsyslog.d/22-loggly.conf
sudo echo '' >> /etc/rsyslog.d/22-loggly.conf
sudo echo '# Send messages to Loggly over TCP using the template.' >> /etc/rsyslog.d/22-loggly.conf
sudo echo '*.* @@logs-01.loggly.com:514;LogglyFormat' >> /etc/rsyslog.d/22-loggly.conf

sudo service rsyslog restart

EOF

  # TODO make this forward the logs of Docker to Loggly

  lifecycle {
    create_before_destroy = true
  }

  security_groups             = [
    "${aws_security_group.gateway.id}",
    "${aws_security_group.private.id}",
    "${aws_security_group.service.id}",
  ]
}

resource "aws_autoscaling_group" "service" {
  vpc_zone_identifier       = [
    "${aws_subnet.frontend-a.id}",
    "${aws_subnet.frontend-b.id}",
  ]
  name                      = "service"
  max_size                  = 30
  min_size                  = 2
  health_check_type         = "EC2"
  health_check_grace_period = 60
  force_delete              = false
  launch_configuration      = "${aws_launch_configuration.service.name}"
  load_balancers            = [
    "${aws_elb.dashboard.name}",
    "${aws_elb.gateway-http.name}",
  ]
  termination_policies      = [
    "OldestInstance",
    "OldestLaunchConfiguration",
    "ClosestToNextInstanceHour",
  ]

  tag {
    key                 = "Name"
    value               = "service"
    propagate_at_launch = true
  }
}

resource "cloudflare_record" "dashboard" {
  domain  = "${var.cloudflare_domain}"
  name    = "dashboard"
  value   = "${aws_elb.dashboard.dns_name}"
  type    = "CNAME"
  ttl     = 1
  proxied = true
}

resource "cloudflare_record" "gateway-http" {
  domain  = "${var.cloudflare_domain}"
  name    = "gateway-http"
  value   = "${aws_elb.gateway-http.dns_name}"
  type    = "CNAME"
  ttl     = 1
  proxied = true
}
