resource "aws_subnet" "backend-a" {
  vpc_id                  = "${aws_vpc.tapglue.id}"
  map_public_ip_on_launch = false

  cidr_block              = "10.0.24.0/22"
  availability_zone       = "${var.zone-a}"

  tags {
    Name = "Backend A"
  }
}

resource "aws_subnet" "backend-b" {
  vpc_id                  = "${aws_vpc.tapglue.id}"
  map_public_ip_on_launch = false

  cidr_block              = "10.0.28.0/22"
  availability_zone       = "${var.zone-b}"

  tags {
    Name = "Backend B"
  }
}

# Routing tables
resource "aws_route_table_association" "backend-a" {
  subnet_id      = "${aws_subnet.backend-a.id}"
  route_table_id = "${aws_route_table.to-nat.id}"
}

resource "aws_route_table_association" "backend-b" {
  subnet_id      = "${aws_subnet.backend-b.id}"
  route_table_id = "${aws_route_table.to-nat.id}"
}
