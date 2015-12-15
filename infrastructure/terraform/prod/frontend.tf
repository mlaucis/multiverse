resource "aws_subnet" "frontend-a" {
  availability_zone       = "${var.zone-a}"
  cidr_block              = "10.0.12.0/22"
  map_public_ip_on_launch = false
  vpc_id                  = "${aws_vpc.tapglue.id}"

  tags {
    Name = "Frontend A"
  }
}

resource "aws_subnet" "frontend-b" {
  availability_zone       = "${var.zone-b}"
  cidr_block              = "10.0.16.0/22"
  map_public_ip_on_launch = false
  vpc_id                  = "${aws_vpc.tapglue.id}"

  tags {
    Name = "Frontend B"
  }
}

# Routing tables
resource "aws_route_table_association" "frontend-a" {
  subnet_id      = "${aws_subnet.frontend-a.id}"
  route_table_id = "${aws_route_table.to-nat.id}"
}

resource "aws_route_table_association" "frontend-b" {
  subnet_id      = "${aws_subnet.frontend-b.id}"
  route_table_id = "${aws_route_table.to-nat.id}"
}

# ELB
resource "aws_elb" "frontend" {
  connection_draining         = true
  connection_draining_timeout = 10
  cross_zone_load_balancing   = true
  idle_timeout                = 30
  name                        = "frontend-prod"
  subnets                     = [
    "${aws_subnet.public-a.id}",
    "${aws_subnet.public-b.id}",
  ]
  security_groups             = [
    "${aws_security_group.loadbalancer.id}",
  ]

  access_logs {
    bucket = "tapglue-logs"
    interval = 5
  }

  listener {
    instance_port     = 8083
    instance_protocol = "https"

    lb_port           = 443
    lb_protocol       = "https"

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
    Name = "frontend"
  }
}
