resource "aws_security_group" "gateway" {
  description = "Gateway node firewall rules"
  name        = "gateway"
  vpc_id      = "${aws_vpc.tapglue.id}"

  tags {
    Name = "gateway"
  }
}

resource "aws_security_group_rule" "gateway_http_in" {
  from_port                 = 80
  to_port                   = 80
  type                      = "ingress"
  protocol                  = "tcp"
  security_group_id         = "${aws_security_group.gateway.id}"
  source_security_group_id  = "${aws_security_group.loadbalancer.id}"
}

resource "aws_security_group_rule" "gateway_https_in" {
  from_port                 = 443
  to_port                   = 443
  type                      = "ingress"
  protocol                  = "tcp"
  security_group_id         = "${aws_security_group.gateway.id}"
  source_security_group_id  = "${aws_security_group.loadbalancer.id}"
}

resource "aws_security_group_rule" "gateway_service_in" {
  from_port                 = 8083
  to_port                   = 8083
  type                      = "ingress"
  protocol                  = "tcp"
  security_group_id         = "${aws_security_group.gateway.id}"
  source_security_group_id  = "${aws_security_group.loadbalancer.id}"
}

resource "aws_security_group" "loadbalancer" {
  description = "Loadbalancer firewall rules"
  name        = "loadbalancer"
  vpc_id      = "${aws_vpc.tapglue.id}"

  tags {
    Name = "loadbalancer"
  }
}

resource "aws_security_group_rule" "loadbalancer_cloudflae_http" {
  from_port         = 80
  to_port           = 80
  type              = "ingress"
  protocol          = "tcp"
  security_group_id = "${aws_security_group.loadbalancer.id}"
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

resource "aws_security_group_rule" "loadbalancer_cloudflae_https" {
  from_port         = 443
  to_port           = 443
  type              = "ingress"
  protocol          = "tcp"
  security_group_id = "${aws_security_group.loadbalancer.id}"
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

resource "aws_security_group_rule" "loadbalancer_http_out" {
  from_port                 = 80
  to_port                   = 80
  type                      = "egress"
  protocol                  = "tcp"
  security_group_id         = "${aws_security_group.loadbalancer.id}"
  source_security_group_id  = "${aws_security_group.gateway.id}"
}

resource "aws_security_group_rule" "loadbalancer_https_out" {
  from_port                 = 443
  to_port                   = 443
  type                      = "egress"
  protocol                  = "tcp"
  security_group_id         = "${aws_security_group.loadbalancer.id}"
  source_security_group_id  = "${aws_security_group.gateway.id}"
}

resource "aws_security_group_rule" "loadbalancer_service_out" {
  from_port                 = 8083
  to_port                   = 8083
  type                      = "egress"
  protocol                  = "tcp"
  security_group_id         = "${aws_security_group.loadbalancer.id}"
  source_security_group_id  = "${aws_security_group.gateway.id}"
}

resource "aws_security_group" "nat" {
  description = "NAT node firewall rules"
  name        = "nat"
  vpc_id      = "${aws_vpc.tapglue.id}"

  tags {
    Name = "nat"
  }
}

resource "aws_security_group_rule" "nat_http_in" {
  from_port                 = 80
  to_port                   = 80
  type                      = "ingress"
  protocol                  = "tcp"
  security_group_id         = "${aws_security_group.nat.id}"
  source_security_group_id  = "${aws_security_group.private.id}"
}

resource "aws_security_group_rule" "nat_http_out" {
  from_port         = 80
  to_port           = 80
  type              = "egress"
  protocol          = "tcp"
  security_group_id = "${aws_security_group.nat.id}"
  cidr_blocks = [
    "0.0.0.0/0"
  ]
}
resource "aws_security_group_rule" "nat_https_in" {
  from_port                 = 443
  to_port                   = 443
  type                      = "ingress"
  protocol                  = "tcp"
  security_group_id         = "${aws_security_group.nat.id}"
  source_security_group_id  = "${aws_security_group.private.id}"
}

resource "aws_security_group_rule" "nat_https_out" {
  from_port         = 443
  to_port           = 443
  type              = "egress"
  protocol          = "tcp"
  security_group_id = "${aws_security_group.nat.id}"
  cidr_blocks = [
    "0.0.0.0/0"
  ]
}

resource "aws_security_group_rule" "nat_ssh_in" {
  from_port         = 22
  to_port           = 22
  type              = "ingress"
  protocol          = "tcp"
  security_group_id = "${aws_security_group.nat.id}"
  cidr_blocks       = [
    "0.0.0.0/0",
  ]
}

resource "aws_security_group_rule" "nat_ssh_out" {
  from_port                 = 22
  to_port                   = 22
  type                      = "egress"
  protocol                  = "tcp"
  security_group_id         = "${aws_security_group.nat.id}"
  source_security_group_id  = "${aws_security_group.private.id}"
}

resource "aws_security_group" "platform" {
  description = "Platform node firewall rules"
  name        = "platform"
  vpc_id      = "${aws_vpc.tapglue.id}"

  tags {
    Name = "platform"
  }
}

resource "aws_security_group_rule" "platform_metrics_out" {
  from_port                 = 9000
  to_port                   = 9100
  type                      = "egress"
  protocol                  = "tcp"
  security_group_id         = "${aws_security_group.platform.id}"
  source_security_group_id  = "${aws_security_group.service.id}"
}

resource "aws_security_group_rule" "platform_mysql_in" {
  from_port                 = 3306
  to_port                   = 3306
  type                      = "ingress"
  protocol                  = "tcp"
  security_group_id         = "${aws_security_group.platform.id}"
  source_security_group_id  = "${aws_security_group.service.id}"
}

resource "aws_security_group_rule" "platform_postgres_in" {
  from_port                 = 5432
  to_port                   = 5432
  type                      = "ingress"
  protocol                  = "tcp"
  security_group_id         = "${aws_security_group.platform.id}"
  source_security_group_id  = "${aws_security_group.service.id}"
}

resource "aws_security_group_rule" "platform_redis_in" {
  from_port                 = 6379
  to_port                   = 6379
  type                      = "ingress"
  protocol                  = "tcp"
  security_group_id         = "${aws_security_group.platform.id}"
  source_security_group_id  = "${aws_security_group.service.id}"
}

resource "aws_security_group" "private" {
  description = "Private node firewall rules"
  name        = "private"
  vpc_id      = "${aws_vpc.tapglue.id}"

  tags {
    Name = "private"
  }
}

resource "aws_security_group_rule" "private_http_out" {
  from_port                 = 80
  to_port                   = 80
  type                      = "egress"
  protocol                  = "tcp"
  security_group_id         = "${aws_security_group.private.id}"
  source_security_group_id  = "${aws_security_group.nat.id}"
}

resource "aws_security_group_rule" "private_https_out" {
  from_port                 = 443
  to_port                   = 443
  type                      = "egress"
  protocol                  = "tcp"
  security_group_id         = "${aws_security_group.private.id}"
  source_security_group_id  = "${aws_security_group.nat.id}"
}

resource "aws_security_group_rule" "private_ssh_in" {
  from_port                 = 22
  to_port                   = 22
  type                      = "ingress"
  protocol                  = "tcp"
  security_group_id         = "${aws_security_group.private.id}"
  source_security_group_id  = "${aws_security_group.nat.id}"
}

resource "aws_security_group" "service" {
  description = "Service node firewall rules"
  name        = "service"
  vpc_id      = "${aws_vpc.tapglue.id}"

  tags {
    Name = "service"
  }
}

resource "aws_security_group_rule" "service_metrics_in" {
  from_port                 = 9000
  to_port                   = 9100
  type                      = "ingress"
  protocol                  = "tcp"
  security_group_id         = "${aws_security_group.service.id}"
  source_security_group_id  = "${aws_security_group.platform.id}"
}

resource "aws_security_group_rule" "service_mysql_out" {
  from_port                 = 3306
  to_port                   = 3306
  type                      = "egress"
  protocol                  = "tcp"
  security_group_id         = "${aws_security_group.service.id}"
  source_security_group_id  = "${aws_security_group.platform.id}"
}

resource "aws_security_group_rule" "service_postgres_out" {
  from_port                 = 5432
  to_port                   = 5432
  type                      = "egress"
  protocol                  = "tcp"
  security_group_id         = "${aws_security_group.service.id}"
  source_security_group_id  = "${aws_security_group.platform.id}"
}

resource "aws_security_group_rule" "service_redis_out" {
  from_port                 = 6379
  to_port                   = 6379
  type                      = "egress"
  protocol                  = "tcp"
  security_group_id         = "${aws_security_group.service.id}"
  source_security_group_id  = "${aws_security_group.platform.id}"
}
