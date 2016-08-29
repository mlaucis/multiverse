resource "aws_instance" "bastion" {
  ami                     = "${var.ami_bastion}"
  availability_zone       = "${var.zone-bastion}"
  instance_type           = "${var.bastion-size}"
  vpc_security_group_ids  = [
    "${aws_security_group.nat.id}"]
  subnet_id               = "${aws_subnet.public-a.id}"
  tags {
    Name = "Bastion Host"
  }
}

resource "aws_eip" "bastion" {
  instance = "${aws_instance.bastion.id}"
  vpc      = true
}
