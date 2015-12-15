resource "aws_subnet" "corporate-a" {
  vpc_id                  = "${aws_vpc.tapglue.id}"
  map_public_ip_on_launch = false

  cidr_block              = "10.0.44.0/24"
  availability_zone       = "${var.zone-a}"

  tags {
    Name = "Corporate A"
  }
}

resource "aws_subnet" "corporate-b" {
  vpc_id                  = "${aws_vpc.tapglue.id}"
  map_public_ip_on_launch = false

  cidr_block              = "10.0.45.0/24"
  availability_zone       = "${var.zone-b}"

  tags {
    Name = "Corporate B"
  }
}

# Routing tables
resource "aws_route_table_association" "corporate-a" {
  subnet_id      = "${aws_subnet.corporate-a.id}"
  route_table_id = "${aws_route_table.to-nat.id}"
}

resource "aws_route_table_association" "corporate-b" {
  subnet_id      = "${aws_subnet.corporate-b.id}"
  route_table_id = "${aws_route_table.to-nat.id}"
}

# ELB
resource "aws_elb" "corporate" {
  name                        = "corporate-prod"
  cross_zone_load_balancing   = true
  idle_timeout                = 300
  connection_draining         = true
  connection_draining_timeout = 10
  subnets                     = [
    "${aws_subnet.public-a.id}",
    "${aws_subnet.public-b.id}"]
  security_groups             = [
    "${aws_security_group.loadbalancer.id}",
  ]

  access_logs {
    bucket = "tapglue-logs"
    interval = 5
  }

  listener {
    lb_port           = 80
    lb_protocol       = "http"

    instance_port     = 80
    instance_protocol = "http"
  }

  listener {
    lb_port           = 443
    lb_protocol       = "https"

    instance_port     = 443
    instance_protocol = "https"

    ssl_certificate_id = "${aws_iam_server_certificate.self-signed.arn}"
  }

  health_check {
    healthy_threshold   = 2
    unhealthy_threshold = 2
    timeout             = 2
    target              = "HTTP:80/"
    interval            = 5
  }

  tags {
    Name = "Corporate"
  }
}
