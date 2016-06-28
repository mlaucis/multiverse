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

resource "aws_ecr_repository" "pganalyze" {
  provider = "aws.us-east-1"
  name     = "pganalyze"
}

resource "aws_ecr_repository_policy" "pganalyze-deployment" {
  provider = "aws.us-east-1"
  repository = "${aws_ecr_repository.pganalyze.name}"
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

resource "aws_ecr_repository" "reporter" {
  provider = "aws.us-east-1"
  name     = "reporter"
}

resource "aws_ecr_repository_policy" "reporter-deployment" {
  provider = "aws.us-east-1"
  repository = "${aws_ecr_repository.reporter.name}"
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

resource "aws_ecr_repository" "sims" {
  provider = "aws.us-east-1"
  name     = "sims"
}

resource "aws_ecr_repository_policy" "sims-deployment" {
  provider = "aws.us-east-1"
  repository = "${aws_ecr_repository.sims.name}"
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
    target              = "HTTPS:8083/health-45016490610398192"
    timeout             = 2
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
    "image": "775034650473.dkr.ecr.us-east-1.amazonaws.com/dashboard:${var.version.dashboard}",
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
  cluster                             = "${aws_ecs_cluster.service.id}"
  depends_on                          = [
    "aws_iam_role_policy.ecsELB",
  ]
  deployment_maximum_percent          = 200
  deployment_minimum_healthy_percent  = 50
  desired_count                       = 2
  iam_role                            = "${aws_iam_role.ecsELB.arn}"
  name                                = "dashboard"
  task_definition                     = "${aws_ecs_task_definition.dashboard.arn}"

  load_balancer {
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
    "command": [
      "./gateway-http",
      "-aws.id", "${aws_iam_access_key.state-change-sr.id}",
      "-aws.secret", "${aws_iam_access_key.state-change-sr.secret}",
      "-aws.region", "${var.vpc-region}",
      "-source", "sqs"
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

resource "aws_ecs_service" "gateway-http" {
  cluster                             = "${aws_ecs_cluster.service.id}"
  depends_on                          = [
    "aws_iam_role_policy.ecsELB",
  ]
  deployment_maximum_percent          = 200
  deployment_minimum_healthy_percent  = 50
  desired_count                       = 2
  iam_role                            = "${aws_iam_role.ecsELB.arn}"
  name                                = "gateway-http"
  task_definition                     = "${aws_ecs_task_definition.gateway-http.arn}"

  load_balancer {
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

resource "aws_ecs_task_definition" "pganalyze" {
  family                = "pganalyze"
  container_definitions = <<EOF
[
  {
    "command": [
    "/usr/bin/python",
    "./pganalyze-collector.zip",
    "--config",
    "/.pganalyze_collector.conf"
    ],
    "cpu": 512,
    "dnsSearchDomains": [
      "${var.env}.${var.region}"
    ],
    "essential": true,
    "image": "775034650473.dkr.ecr.us-east-1.amazonaws.com/pganalyze:${var.version.pganalyze}",
    "logConfiguration": {
      "logDriver": "syslog"
    },
    "memory": 1024,
    "name": "pganalyze",
    "portMappings": [],
    "readonlyRootFilesystem": true,
    "workingDirectory": "/"
  }
]
EOF
}

resource "aws_ecs_task_definition" "reporter" {
  family                = "reporter"
  container_definitions = <<EOF
[
  {
    "command": [
      "./reporter",
      "-pg.url", "postgres://${var.rds_username}:${var.rds_password}@db-master.service:5432/${var.rds_db_name}?sslmode=disable&connect_timeout=5",
      "-slack.channel", "reports",
      "-slack.token", "${var.slack_token}"
    ],
    "cpu": 256,
    "dnsSearchDomains": [
      "${var.env}.${var.region}"
    ],
    "essential": true,
    "image": "775034650473.dkr.ecr.us-east-1.amazonaws.com/reporter:${var.version.reporter}",
    "logConfiguration": {
      "logDriver": "syslog"
    },
    "memory": 512,
    "name": "reporter",
    "portMappings": [],
    "readonlyRootFilesystem": true,
    "workingDirectory": "/tapglue/"
  }
]
EOF
}

resource "aws_ecs_task_definition" "sims" {
  family                = "sims"
  container_definitions = <<EOF
[
  {
    "command": [
      "./sims",
      "-aws.id", "${aws_iam_access_key.state-change-sr.id}",
      "-aws.secret", "${aws_iam_access_key.state-change-sr.secret}",
      "-aws.region", "${var.vpc-region}",
      "-postgres.url", "postgres://${var.rds_username}:${var.rds_password}@db-master.service:5432/${var.rds_db_name}?connect_timeout=5"
    ],
    "cpu": 256,
    "dnsSearchDomains": [
      "${var.env}.${var.region}"
    ],
    "essential": true,
    "image": "775034650473.dkr.ecr.us-east-1.amazonaws.com/sims:${var.version.sims}",
    "logConfiguration": {
      "logDriver": "syslog"
    },
    "memory": 512,
    "name": "gateway-http",
    "portMappings": [
      {
        "containerPort": 9001,
        "hostPort": 9001
      }
    ],
    "readonlyRootFilesystem": true,
    "workingDirectory": "/tapglue/"
  }
]
EOF
}

resource "aws_ecs_service" "sims" {
  cluster                             = "${aws_ecs_cluster.service.id}"
  deployment_maximum_percent          = 200
  deployment_minimum_healthy_percent  = 50
  desired_count                       = 2
  name                                = "sims"
  task_definition                     = "${aws_ecs_task_definition.sims.arn}"
}

resource "aws_launch_configuration" "service" {
  associate_public_ip_address = false
  ebs_optimized               = false
  enable_monitoring           = true
  iam_instance_profile        = "${aws_iam_instance_profile.service.name}"
  image_id                    = "${var.ami_container}"
  instance_type               = "m4.large"
  key_name                    = "${aws_key_pair.debug.key_name}"
  user_data                   = <<EOF
#!/bin/bash
echo ECS_CLUSTER=service >> /etc/ecs/ecs.config

# Install loggly security credentials
mkdir -pv /etc/rsyslog.d/keys/ca.d
cd /etc/rsyslog.d/keys/ca.d/
curl -O https://logdog.loggly.com/media/logs-01.loggly.com_sha12.crt

# Rsyslog for Loggly

sudo yum install -y rsyslog-gnutls
sudo mkdir -p /var/spool/rsyslog

echo '$template LogglyFormat,"<%pri%>%protocol-version% %timestamp:::date-rfc3339% %HOSTNAME% %app-name% %procid% %msgid% [d2e7097f-25aa-497a-a9e3-d691bd4ec7ab@41058 tag=\"service.prod.eu-central-1\"] %msg%\n"

# Setup disk assisted queues
$WorkDirectory /var/spool/rsyslog # where to place spool files
$ActionQueueFileName fwdRule1     # unique name prefix for spool files
$ActionQueueMaxDiskSpace 100m     # 1gb space limit (use as much as possible)
$ActionQueueSaveOnShutdown on     # save messages to disk on shutdown
$ActionQueueType LinkedList       # run asynchronously
$ActionResumeRetryCount -1        # infinite retries if host is down

# RsyslogGnuTLS
$DefaultNetstreamDriverCAFile /etc/rsyslog.d/keys/ca.d/logs-01.loggly.com_sha12.crt
$ActionSendStreamDriver gtls
$ActionSendStreamDriverMode 1
$ActionSendStreamDriverAuthMode x509/name
$ActionSendStreamDriverPermittedPeer *.loggly.com
*.* @@logs-01.loggly.com:6514;LogglyFormat
' | sudo tee /etc/rsyslog.d/22-loggly.conf > /dev/null

sudo service rsyslog restart

echo '#!/bin/sh

/usr/sbin/logrotate /etc/logrotate.hourly.conf >/dev/null 2>&1
EXITVALUE=$?
if [ $EXITVALUE != 0 ]; then
    /usr/bin/logger -t logrotate "ALERT exited abnormally with [$EXITVALUE]"
fi
exit 0
' | sudo tee /etc/cron.hourly/logrotate > /dev/null

sudo chmod +x /etc/cron.hourly/logrotate

echo '/var/log/messages {
    compress
    create
    daily
    rotate 5
    size 100M
    postrotate
	/bin/kill -HUP `cat /var/run/syslogd.pid 2> /dev/null` 2> /dev/null || true
    endscript
}' | sudo tee /etc/logrotate.hourly.conf > /dev/null

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

resource "aws_key_pair" "debug" {
  key_name    = "debug"
  public_key  = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQCuFsJxH52k7iI4mseWljlbQhwIfbpVPuDCTOBo6YtI7xL3f3jfme4fqziwt+iqavRW2MgGsgoYGITNYstZa5zzT4Zo6CTZ0XpeLYZrrXQOxXrXjesRA478bCsU4gpCrPiy5Uzw3e2d1HLF/deLjnmREshzqaEQKoL8tzG51esBTIna+M5aWD0AGPFotO3J2sFTRnbAIxeVj4bKWAfaE2+WG1MX1VemDGeGrHmW6UbPoymHOD7Y5c/F00Bv+Pgk5LwCyRCvEzMLbl2GHpEJd3vcouwEToyADlN1rXc+85SfVtlwS8F3fX6vqjQ/2fMzG4syaDEeUJLsBcE2glNIwDH/ debug"
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
