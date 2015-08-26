resource "aws_route_table" "to-nat" {
  vpc_id = "${aws_vpc.prod.id}"

  route {
    cidr_block  = "0.0.0.0/0"
    instance_id = "${aws_instance.nat.id}"
  }

  tags {
    Name = "to-nat"
  }
}

# Security groups
resource "aws_security_group" "from-nat" {
  vpc_id      = "${aws_vpc.prod.id}"
  name        = "from-nat"
  description = "Allow services from the private subnet through NAT"

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = [
      "0.0.0.0/0"]
  }

  ingress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = [
      "${aws_subnet.corporate-a.cidr_block}",
      "${aws_subnet.corporate-b.cidr_block}",
      "${aws_subnet.frontend-a.cidr_block}",
      "${aws_subnet.frontend-b.cidr_block}",
      "${aws_subnet.backend-a.cidr_block}",
      "${aws_subnet.backend-b.cidr_block}"]
  }

  tags {
    Name = "From NAT"
  }
}

resource "aws_security_group" "to-nat" {
  vpc_id      = "${aws_vpc.prod.id}"
  name        = "to-nat"
  description = "Allow services from the private subnet through NAT"

  ingress {
    from_port       = 0
    to_port         = 0
    protocol        = "-1"
    security_groups = [
      "${aws_security_group.from-nat.id}"]
  }

  egress {
    from_port       = 0
    to_port         = 0
    protocol        = "-1"
    security_groups = [
      "${aws_security_group.from-nat.id}"]
  }

  tags {
    Name = "To NAT"
  }
}

# Instance
resource "aws_instance" "nat" {
  ami                         = "${var.ami_nat}"
  availability_zone           = "${var.zone-nat}"
  instance_type               = "${var.nat-size}"
  security_groups             = [
    "${aws_security_group.from-nat.id}"]
  subnet_id                   = "${aws_subnet.public-a.id}"
  associate_public_ip_address = true
  source_dest_check           = false

  tags {
    Name = "NAT instance"
  }
}
