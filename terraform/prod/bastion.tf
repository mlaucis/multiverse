# Bastion instance
resource "aws_security_group" "bastion" {
  vpc_id      = "${aws_vpc.prod.id}"
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
      "${aws_subnet.frontend-a.cidr_block}"]
  }
  egress {
    from_port   = 22
    to_port     = 22
    protocol    = "tcp"
    cidr_blocks = [
      "${aws_subnet.frontend-a.cidr_block}"]
  }

  egress {
    from_port   = 22
    to_port     = 22
    protocol    = "tcp"
    cidr_blocks = [
      "${aws_subnet.backend-a.cidr_block}"]
  }
  egress {
    from_port   = 22
    to_port     = 22
    protocol    = "tcp"
    cidr_blocks = [
      "${aws_subnet.backend-a.cidr_block}"]
  }

  tags {
    Name = "Bastion"
  }
}

resource "aws_instance" "bastion" {
  ami               = "${var.aws_ubuntu_ami}"
  availability_zone = "${var.zone-bastion}"
  instance_type     = "${var.bastion-size}"
  key_name          = "${var.aws_key_name}"
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
