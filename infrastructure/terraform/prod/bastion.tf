resource "aws_security_group" "bastion" {
  vpc_id      = "${aws_vpc.tapglue.id}"
  name        = "bastion"
  description = "Allow SSH traffic from the internet"

  ingress {
    from_port   = 22
    to_port     = 22
    protocol    = "tcp"
    cidr_blocks = [
      "0.0.0.0/0"]
  }

  egress {
    from_port   = 22
    to_port     = 22
    protocol    = "tcp"
    cidr_blocks = [
      "${aws_subnet.backend-a.cidr_block}",
      "${aws_subnet.backend-b.cidr_block}",
      "${aws_subnet.corporate-a.cidr_block}",
      "${aws_subnet.corporate-b.cidr_block}",
      "${aws_subnet.frontend-a.cidr_block}",
      "${aws_subnet.frontend-b.cidr_block}",
      "${aws_subnet.monitoring-a.cidr_block}",
      "${aws_subnet.monitoring-b.cidr_block}",
      "${aws_subnet.public-a.cidr_block}",
      "${aws_subnet.public-b.cidr_block}",
    ]
  }

  tags {
    Name = "Bastion"
  }
}

resource "aws_instance" "bastion" {
  ami               = "${var.ami_bastion}"
  availability_zone = "${var.zone-bastion}"
  instance_type     = "${var.bastion-size}"
  security_groups   = [
    "${aws_security_group.bastion.id}"]
  subnet_id         = "${aws_subnet.public-a.id}"
  tags {
    Name = "Bastion Host"
  }
}

resource "aws_eip" "bastion" {
  instance = "${aws_instance.bastion.id}"
  vpc      = true
}
