resource "aws_ecr_repository" "intaker" {
  provider = "aws.us-east-1"
  name     = "intaker"
}

resource "aws_ecr_repository" "corporate" {
  provider = "aws.us-east-1"
  name     = "corporate"
}

resource "aws_ecr_repository_policy" "deployment" {
  repository = "${aws_ecr_repository.intaker.name}"
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
                    "arn:aws:iam::775034650473:role/ecsServiceRole",
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

resource "aws_ecs_cluster" "production-intaker" {
  name = "production-intaker"
}

resource "aws_iam_role_policy" "TapglueEC2ContainerServiceforEC2Role" {
  name   = "TapglueEC2ContainerServiceforEC2Role"
  role   = "${aws_iam_role.tapglueEcsInstanceRole.id}"
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

resource "aws_iam_role_policy" "tapglueEcsInstanceRolePolicy" {
  name   = "tapglueEcsInstanceRolePolicy"
  role   = "${aws_iam_role.tapglueEcsInstanceRole.id}"
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
        "ecs:Submit*"
      ],
      "Resource": [
        "*"
      ]
    }
  ]
}
EOF
}

resource "aws_iam_role" "tapglueEcsInstanceRole" {
  name               = "tapglueEcsInstanceRole"
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

resource "aws_iam_role_policy" "tapglueEcsElbRolePolicy" {
  name   = "tapglueEcsElbRole"
  role   = "${aws_iam_role.tapglueEcsElbRole.id}"
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

resource "aws_iam_role" "tapglueEcsElbRole" {
  name               = "tapglueEcsElbRole"
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

resource "aws_elb" "container-frontend" {
  connection_draining         = true
  connection_draining_timeout = 10
  cross_zone_load_balancing   = true
  idle_timeout                = 30
  name                        = "container-frontend-prod"
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
    Name = "container-frontend"
  }
}

resource "aws_elb" "container-corporate" {
  connection_draining         = true
  connection_draining_timeout = 10
  cross_zone_load_balancing   = true
  idle_timeout                = 30
  name                        = "container-corporate-prod"
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
    target              = "HTTPS:8083/health-45016490610398192"
    unhealthy_threshold = 2
  }

  tags {
    Name = "container-corporate"
  }
}

resource "aws_ecs_task_definition" "intaker" {
  family                = "intaker"
  container_definitions = <<EOF
[
  {
    "name": "intaker",
    "image": "775034650473.dkr.ecr.us-east-1.amazonaws.com/intaker:1632",
    "cpu": 1024,
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

resource "aws_ecs_task_definition" "corporate" {
  family                = "corporate"
  container_definitions = <<EOF
[
  {
    "name": "corporate",
    "image": "775034650473.dkr.ecr.us-east-1.amazonaws.com/corporate:1632",
    "cpu": 512,
    "memory": 512,
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

resource "aws_ecs_service" "intaker" {
  name            = "intaker"
  cluster         = "${aws_ecs_cluster.production-intaker.id}"
  task_definition = "${aws_ecs_task_definition.intaker.arn}"
  desired_count   = 1
  iam_role = "${aws_iam_role.tapglueEcsElbRole.arn}"
  depends_on = [
    "aws_iam_role_policy.tapglueEcsElbRolePolicy"]

  load_balancer = {
    elb_name = "${aws_elb.container-frontend.id}"
    container_name = "intaker"
    container_port = 8083
  }
}

resource "aws_ecs_service" "corporate" {
  name            = "corporate"
  cluster         = "${aws_ecs_cluster.production-intaker.id}"
  task_definition = "${aws_ecs_task_definition.corporate.arn}"
  desired_count   = 1
  iam_role = "${aws_iam_role.tapglueEcsElbRole.arn}"
  depends_on = [
    "aws_iam_role_policy.tapglueEcsElbRolePolicy"]

  load_balancer = {
    elb_name = "${aws_elb.container-corporate.id}"
    container_name = "corporate"
    container_port = 443
  }
}

resource "aws_iam_instance_profile" "container-frontend" {
  name  = "container-frontend"
  roles = [
    "${aws_iam_role.tapglueEcsInstanceRole.name}"]
}

resource "aws_launch_configuration" "container-frontend" {
  image_id                    = "${var.ami_container}"
  instance_type               = "t2.medium"
  associate_public_ip_address = false
  enable_monitoring           = true
  ebs_optimized               = false
  iam_instance_profile        = "${aws_iam_instance_profile.container-frontend.name}"

  user_data                   = <<EOF
#!/bin/bash
echo ECS_CLUSTER=production-intaker >> /etc/ecs/ecs.config

sudo echo '$WorkDirectory /var/spool/rsyslog # where to place spool files' > /etc/rsyslog.d/22-loggly.conf
sudo echo '$ActionQueueFileName fwdRule1     # unique name prefix for spool files' >> /etc/rsyslog.d/22-loggly.conf
sudo echo '$ActionQueueMaxDiskSpace 100m     # 1gb space limit (use as much as possible)' >> /etc/rsyslog.d/22-loggly.conf
sudo echo '$ActionQueueSaveOnShutdown on     # save messages to disk on shutdown' >> /etc/rsyslog.d/22-loggly.conf
sudo echo '$ActionQueueType LinkedList       # run asynchronously' >> /etc/rsyslog.d/22-loggly.conf
sudo echo '$ActionResumeRetryCount -1        # infinite retries if host is down' >> /etc/rsyslog.d/22-loggly.conf
sudo echo '' >> /etc/rsyslog.d/22-loggly.conf
sudo echo '$template LogglyFormat,"<%pri%>%protocol-version% %timestamp:::date-rfc3339% %HOSTNAME% %app-name% %procid% %msgid% [d2e7097f-25aa-497a-a9e3-d691bd4ec7ab@41058 tag=\"production-intaker\"] %msg%\n"' >> /etc/rsyslog.d/22-loggly.conf
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

# Group
resource "aws_autoscaling_group" "container-frontend" {
  vpc_zone_identifier       = [
    "${aws_subnet.frontend-a.id}",
    "${aws_subnet.frontend-b.id}"]
  name                      = "container-frontend"
  max_size                  = 30
  min_size                  = 1
  health_check_type         = "EC2"
  health_check_grace_period = 60
  force_delete              = false
  launch_configuration      = "${aws_launch_configuration.container-frontend.name}"
  load_balancers            = [
    "${aws_elb.container-frontend.name}",
    "${aws_elb.container-corporate.name}"]
  termination_policies      = [
    "OldestInstance",
    "OldestLaunchConfiguration",
    "ClosestToNextInstanceHour"]

  tag {
    key                 = "Name"
    value               = "ecs-frontend"
    propagate_at_launch = true
  }
}

resource "cloudflare_record" "container-api" {
  domain = "${var.cloudflare_domain}"
  name   = "container-api"
  value  = "${aws_elb.container-frontend.dns_name}"
  type   = "CNAME"
  ttl    = 1
}

resource "cloudflare_record" "container-corporate" {
  domain = "${var.cloudflare_domain}"
  name   = "container-corporate"
  value  = "${aws_elb.container-corporate.dns_name}"
  type   = "CNAME"
  ttl    = 1
}
