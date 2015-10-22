resource "aws_route_table" "to-nat" {
  vpc_id = "${aws_vpc.tapglue.id}"

  route {
    cidr_block  = "0.0.0.0/0"
    instance_id = "${aws_instance.nat.id}"
  }

  tags {
    Name = "nat"
  }
}

# Instance
resource "aws_instance" "nat" {
  ami                         = "${var.ami_nat}"
  associate_public_ip_address = true
  availability_zone           = "${var.zone-nat}"
  instance_type               = "${var.nat-size}"
  subnet_id                   = "${aws_subnet.public-a.id}"
  source_dest_check           = false

  security_groups             = [
    "${aws_security_group.nat.id}",
  ]

  tags {
    Name = "nat"
  }
}
