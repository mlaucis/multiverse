resource "aws_subnet" "frontend-a" {
  vpc_id                  = "${aws_vpc.tapglue.id}"
  map_public_ip_on_launch = false

  cidr_block              = "10.0.12.0/22"
  availability_zone       = "${var.zone-a}"

  tags {
    Name = "Frontend A"
  }
}

resource "aws_subnet" "frontend-b" {
  vpc_id                  = "${aws_vpc.tapglue.id}"
  map_public_ip_on_launch = false

  cidr_block              = "10.0.16.0/22"
  availability_zone       = "${var.zone-b}"

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

# Security groups
resource "aws_security_group" "frontend-elb-inet" {
  vpc_id      = "${aws_vpc.tapglue.id}"
  name        = "frontend-elb-inet"
  description = "Allow Internet traffic to and from ELB"

  ingress {
    from_port   = 443
    to_port     = 443
    protocol    = "tcp"
    cidr_blocks = [
      "0.0.0.0/0"]
  }

  egress {
    from_port   = 443
    to_port     = 443
    protocol    = "tcp"
    cidr_blocks = [
      "0.0.0.0/0"]
  }

  tags {
    Name = "Internet to ELB Frontend"
  }
}

resource "aws_security_group" "frontend-elb-vpc" {
  vpc_id      = "${aws_vpc.tapglue.id}"
  name        = "frontend-elb"
  description = "Allow EC2 traffic to and from the ELB"

  ingress {
    from_port   = 8083
    to_port     = 8083
    protocol    = "tcp"
    cidr_blocks = [
      "0.0.0.0/0"]
  }

  egress {
    from_port   = 8083
    to_port     = 8083
    protocol    = "tcp"
    cidr_blocks = [
      "0.0.0.0/0"]
  }

  tags {
    Name = "ELB Frontend to EC2"
  }
}

resource "aws_security_group" "frontend-elb-ec2" {
  vpc_id      = "${aws_vpc.tapglue.id}"
  name        = "frontend-elb-ec2"
  description = "Allow Traffic from ELB to EC2"

  ingress {
    from_port = 8083
    to_port   = 8083
    protocol  = "tcp"
    self      = true
  }

  egress {
    from_port       = 8083
    to_port         = 8083
    protocol        = "tcp"
    security_groups = [
      "${aws_security_group.frontend-elb-vpc.id}"]
  }

  tags {
    Name = "Allow Traffic from ELB to EC2"
  }
}

resource "aws_security_group" "frontend-ssh" {
  vpc_id      = "${aws_vpc.tapglue.id}"
  name        = "frontend-ssh"
  description = "Allow SSH traffic from the Bastion host"

  ingress {
    from_port       = 22
    to_port         = 22
    protocol        = "tcp"
    security_groups = [
      "${aws_security_group.bastion.id}"]
  }

  egress {
    from_port       = 22
    to_port         = 22
    protocol        = "tcp"
    security_groups = [
      "${aws_security_group.bastion.id}"]
  }

  tags {
    Name = "SSH from Bastion to Frontend"
  }
}

# ELB
resource "aws_elb" "frontend" {
  name                        = "frontend"
  cross_zone_load_balancing   = true
  idle_timeout                = 300
  connection_draining         = true
  connection_draining_timeout = 10
  subnets                     = [
    "${aws_subnet.public-a.id}",
    "${aws_subnet.public-b.id}"]
  security_groups             = [
    "${aws_security_group.frontend-elb-inet.id}",
    "${aws_security_group.frontend-elb-ec2.id}"]

  listener {
    lb_port           = 443
    lb_protocol       = "https"

    instance_port     = 8083
    instance_protocol = "https"

    ssl_certificate_id = "${aws_iam_server_certificate.self-signed.arn}"
  }

  health_check {
    healthy_threshold   = 2
    unhealthy_threshold = 2
    timeout             = 2
    target              = "HTTPS:8083/health-45016490610398192"
    interval            = 5
  }

  tags {
    Name = "Frontend"
  }
}
