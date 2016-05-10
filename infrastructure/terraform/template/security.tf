resource "aws_iam_role" "ecs-agent" {
  name                = "ecs-agent"
  assume_role_policy  = <<EOF
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

resource "aws_iam_role_policy" "ecs-agent" {
  name    = "ecs-agent"
  role    = "${aws_iam_role.ecs-agent.id}"
  policy  = <<EOF
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

resource "aws_iam_role_policy" "ecs-sqs" {
  name    = "ecs-sqs"
  role    = "${aws_iam_role.esc-agent.id}"
  policy  = <<EOF
{
   "Version": "2012-10-17",
   "Statement":[{
      "Effect":"Allow",
      "Action": [
        "sqs:SendMessage",
        "sqs:ReceiveMessage",
        "sqs:GetQueueUrl"
      ],
      "Resource":"arn:aws:sqs:*:775034650473:*"
      }
   ]
}
EOF
}


resource "aws_iam_role" "ecs-scheduler" {
  name                = "ecs-scheduler"
  assume_role_policy  = <<EOF
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

resource "aws_iam_role_policy" "ecs-scheduler" {
  name    = "ecs-scheduler"
  role    = "${aws_iam_role.ecs-scheduler.id}"
  policy  = <<EOF
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

resource "aws_iam_instance_profile" "ecs-agent-profile" {
  name  = "ecs-agent-profile"
  roles = [
    "${aws_iam_role.ecs-agent.name}",
  ]
}

resource "aws_security_group" "perimeter" {
  description = "perimeter firewall rules"
  name        = "perimeter"
  vpc_id      = "${aws_vpc.env.id}"

  tags {
    Name = "perimeter"
  }
}

resource "aws_key_pair" "debug" {
  key_name    = "debug"
  public_key  = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQCuFsJxH52k7iI4mseWljlbQhwIfbpVPuDCTOBo6YtI7xL3f3jfme4fqziwt+iqavRW2MgGsgoYGITNYstZa5zzT4Zo6CTZ0XpeLYZrrXQOxXrXjesRA478bCsU4gpCrPiy5Uzw3e2d1HLF/deLjnmREshzqaEQKoL8tzG51esBTIna+M5aWD0AGPFotO3J2sFTRnbAIxeVj4bKWAfaE2+WG1MX1VemDGeGrHmW6UbPoymHOD7Y5c/F00Bv+Pgk5LwCyRCvEzMLbl2GHpEJd3vcouwEToyADlN1rXc+85SfVtlwS8F3fX6vqjQ/2fMzG4syaDEeUJLsBcE2glNIwDH/ debug"
}

resource "aws_security_group_rule" "perimeter_cloudflare_https_in" {
  from_port         = 443
  to_port           = 443
  type              = "ingress"
  protocol          = "tcp"
  security_group_id = "${aws_security_group.perimeter.id}"
  cidr_blocks = [
    "103.21.244.0/22",
    "103.22.200.0/22",
    "103.31.4.0/22",
    "104.16.0.0/12",
    "108.162.192.0/18",
    "141.101.64.0/18",
    "162.158.0.0/15",
    "172.64.0.0/13",
    "173.245.48.0/20",
    "188.114.96.0/20",
    "190.93.240.0/20",
    "197.234.240.0/22",
    "198.41.128.0/17",
    "199.27.128.0/21",
  ]
}

resource "aws_security_group_rule" "perimeter_grafana_out" {
  from_port                 = 3000
  to_port                   = 3000
  type                      = "egress"
  protocol                  = "tcp"
  security_group_id         = "${aws_security_group.perimeter.id}"
  source_security_group_id  = "${aws_security_group.platform.id}"
}

resource "aws_security_group_rule" "perimeter_http_out" {
  cidr_blocks       = [
    "0.0.0.0/0"
  ]
  from_port         = 80
  to_port           = 80
  type              = "egress"
  protocol          = "tcp"
  security_group_id = "${aws_security_group.perimeter.id}"
}

resource "aws_security_group_rule" "perimeter_https_out" {
  cidr_blocks       = [
    "0.0.0.0/0"
  ]
  from_port         = 443
  to_port           = 443
  type              = "egress"
  protocol          = "tcp"
  security_group_id = "${aws_security_group.perimeter.id}"
}

resource "aws_security_group_rule" "perimeter_service_out" {
  from_port                 = 8083
  to_port                   = 8083
  type                      = "egress"
  protocol                  = "tcp"
  security_group_id         = "${aws_security_group.perimeter.id}"
  source_security_group_id  = "${aws_security_group.platform.id}"
}

resource "aws_security_group_rule" "perimeter_ssh_in" {
  cidr_blocks       = [
    "0.0.0.0/0",
  ]
  from_port         = 22
  to_port           = 22
  type              = "ingress"
  protocol          = "tcp"
  security_group_id = "${aws_security_group.perimeter.id}"
}

resource "aws_security_group_rule" "perimeter_ssh_out" {
  cidr_blocks       = [
    "10.0.0.0/16",
  ]
  from_port         = 22
  to_port           = 22
  type              = "egress"
  protocol          = "tcp"
  security_group_id = "${aws_security_group.perimeter.id}"
}

resource "aws_security_group" "platform" {
  description = "platform firewall rules"
  name        = "platform"
  vpc_id      = "${aws_vpc.env.id}"

  tags {
    Name = "platform"
  }
}

resource "aws_security_group_rule" "platform_grafana_in" {
  from_port                 = 3000
  to_port                   = 3000
  type                      = "ingress"
  protocol                  = "tcp"
  security_group_id         = "${aws_security_group.platform.id}"
  source_security_group_id  = "${aws_security_group.perimeter.id}"
}

resource "aws_security_group_rule" "platform_http_out" {
  cidr_blocks       = [
    "0.0.0.0/0"
  ]
  from_port         = 80
  to_port           = 80
  type              = "egress"
  protocol          = "tcp"
  security_group_id = "${aws_security_group.platform.id}"
}

resource "aws_security_group_rule" "platform_https_out" {
  cidr_blocks       = [
    "0.0.0.0/0"
  ]
  from_port         = 443
  to_port           = 443
  type              = "egress"
  protocol          = "tcp"
  security_group_id = "${aws_security_group.platform.id}"
}

resource "aws_security_group_rule" "platform_postgres_in" {
  from_port         = 5432
  to_port           = 5432
  type              = "ingress"
  protocol          = "tcp"
  security_group_id = "${aws_security_group.platform.id}"
  self              = true
}

resource "aws_security_group_rule" "platform_postgres_out" {
  from_port         = 5432
  to_port           = 5432
  type              = "egress"
  protocol          = "tcp"
  security_group_id = "${aws_security_group.platform.id}"
  self              = true
}

resource "aws_security_group_rule" "platform_prometheus_in" {
  from_port         = 9000
  to_port           = 9000
  type              = "ingress"
  protocol          = "tcp"
  security_group_id = "${aws_security_group.platform.id}"
  self              = true
}

resource "aws_security_group_rule" "platform_prometheus_out" {
  from_port         = 9000
  to_port           = 9000
  type              = "egress"
  protocol          = "tcp"
  security_group_id = "${aws_security_group.platform.id}"
  self              = true
}

resource "aws_security_group_rule" "platform_redis_in" {
  from_port         = 6379
  to_port           = 6379
  type              = "ingress"
  protocol          = "tcp"
  security_group_id = "${aws_security_group.platform.id}"
  self              = true
}

resource "aws_security_group_rule" "platform_redis_out" {
  from_port         = 6379
  to_port           = 6379
  type              = "egress"
  protocol          = "tcp"
  security_group_id = "${aws_security_group.platform.id}"
  self              = true
}

resource "aws_security_group_rule" "platform_service_in" {
  from_port                 = 8083
  to_port                   = 8083
  type                      = "ingress"
  protocol                  = "tcp"
  security_group_id         = "${aws_security_group.platform.id}"
  source_security_group_id  = "${aws_security_group.perimeter.id}"
}

resource "aws_security_group_rule" "platform_ssh_in" {
  from_port                 = 22
  to_port                   = 22
  type                      = "ingress"
  protocol                  = "tcp"
  security_group_id         = "${aws_security_group.platform.id}"
  source_security_group_id  = "${aws_security_group.perimeter.id}"
}

resource "aws_security_group_rule" "platform_syslog_out" {
  cidr_blocks       = [
    "0.0.0.0/0"
  ]
  from_port         = 6514
  to_port           = 6514
  type              = "egress"
  protocol          = "udp"
  security_group_id = "${aws_security_group.platform.id}"
}