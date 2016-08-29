resource "aws_subnet" "monitoring-a" {
  availability_zone       = "${var.zone-a}"
  cidr_block              = "10.0.46.0/24"
  map_public_ip_on_launch = false
  vpc_id                  = "${aws_vpc.tapglue.id}"


  tags {
    Name = "monitoring-a"
  }
}

resource "aws_subnet" "monitoring-b" {
  availability_zone       = "${var.zone-b}"
  cidr_block              = "10.0.47.0/24"
  map_public_ip_on_launch = false
  vpc_id                  = "${aws_vpc.tapglue.id}"

  tags {
    Name = "monitoring-b"
  }
}

resource "aws_route_table_association" "monitoring-a" {
  subnet_id      = "${aws_subnet.monitoring-a.id}"
  route_table_id = "${aws_route_table.to-nat.id}"
}

resource "aws_route_table_association" "monitoring-b" {
  subnet_id      = "${aws_subnet.monitoring-b.id}"
  route_table_id = "${aws_route_table.to-nat.id}"
}

resource "aws_instance" "monitoring0" {
  ami           = "${var.monitoring_ami}"
  instance_type = "${var.monitoring_instance_type}"
  subnet_id     = "${aws_subnet.monitoring-a.id}"

  vpc_security_group_ids = [
    "${aws_security_group.platform.id}",
    "${aws_security_group.private.id}",
  ]

  tags {
    Name = "monitoring0"
  }
}

resource "aws_elb" "monitoring" {
  cross_zone_load_balancing   = true
  connection_draining         = true
  connection_draining_timeout = 10
  idle_timeout                = 15
  name                        = "monitoring"

  instances = [
    "${aws_instance.monitoring0.id}",
  ]

  listener = {
    instance_port       = 3000
    instance_protocol   = "http"
    lb_port             = 443
    lb_protocol         = "https"
    ssl_certificate_id  = "${aws_iam_server_certificate.self-signed.arn}"
  }

  security_groups = [
    "${aws_security_group.loadbalancer.id}",
  ]

  subnets = [
    "${aws_subnet.public-a.id}",
    "${aws_subnet.public-b.id}",
  ]

  tags = {
    Name = "monitoring"
  }
}
