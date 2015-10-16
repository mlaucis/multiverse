resource "aws_security_group" "monitoring-collect" {
  description = "Metrics collection from the monitoring hosts"
  name        = "monitoring-collect"
  vpc_id      = "${aws_vpc.tapglue.id}"

  egress {
    from_port = 9000
    to_port   = 9100
    protocol  = "tcp"
    cidr_blocks = [
      "${aws_subnet.backend-a.cidr_block}",
      "${aws_subnet.backend-b.cidr_block}",
      "${aws_subnet.corporate-a.cidr_block}",
      "${aws_subnet.corporate-b.cidr_block}",
      "${aws_subnet.frontend-a.cidr_block}",
      "${aws_subnet.frontend-b.cidr_block}",
      "${aws_subnet.monitoring-a.cidr_block}",
      "${aws_subnet.monitoring-b.cidr_block}",
    ]
  }

  tags {
    Name = "monitoring-collect"
  }
}

resource "aws_security_group" "monitored" {
  description = "Allow metrics collection"
  name        = "monitored"
  vpc_id      = "${aws_vpc.tapglue.id}"

  ingress {
    from_port = 9000
    to_port   = 9100
    protocol  = "tcp"
    cidr_blocks = [
      "${aws_subnet.monitoring-a.cidr_block}",
      "${aws_subnet.monitoring-b.cidr_block}",
    ]
  }

  tags {
    Name = "monitored"
  }
}

resource "aws_security_group" "monitoring-ssh" {
  description = "Allow SSH traffic from the Bastion host"
  name        = "monitoring-ssh"
  vpc_id      = "${aws_vpc.tapglue.id}"

  ingress {
    from_port       = 22
    to_port         = 22
    protocol        = "tcp"
    security_groups = [
      "${aws_security_group.bastion.id}",
    ]
  }

  tags {
    Name = "monitoring-ssh"
  }
}

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

resource "aws_instance" "monitoring0" {
  ami           = "${var.monitoring_ami}"
  instance_type = "${var.monitoring_instance_type}"
  subnet_id     = "${aws_subnet.monitoring-a.id}"

  security_groups = [
    "${aws_security_group.monitoring-ssh.id}",
    "${aws_security_group.to-nat.id}",
  ]

  tags {
    Name = "monitoring0"
  }
}

resource "aws_instance" "monitoring1" {
  ami           = "${var.monitoring_ami}"
  instance_type = "${var.monitoring_instance_type}"
  subnet_id     = "${aws_subnet.monitoring-b.id}"

  security_groups = [
    "${aws_security_group.monitoring-ssh.id}",
    "${aws_security_group.to-nat.id}",
  ]

  tags {
    Name = "monitoring1"
  }
}
