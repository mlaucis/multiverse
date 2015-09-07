# Public instances
resource "aws_subnet" "public-a" {
  vpc_id                  = "${aws_vpc.tapglue.id}"
  map_public_ip_on_launch = true

  cidr_block              = "10.0.0.0/22"
  availability_zone       = "${var.zone-a}"

  tags {
    Name = "Public A"
  }
}

resource "aws_subnet" "public-b" {
  vpc_id                  = "${aws_vpc.tapglue.id}"
  map_public_ip_on_launch = true

  cidr_block              = "10.0.4.0/22"
  availability_zone       = "${var.zone-b}"

  tags {
    Name = "Public B"
  }
}

resource "aws_route_table" "public" {
  vpc_id = "${aws_vpc.tapglue.id}"

  route {
    cidr_block = "0.0.0.0/0"
    gateway_id = "${aws_internet_gateway.default.id}"
  }

  tags {
    Name = "public"
  }
}

resource "aws_route_table_association" "public-a" {
  subnet_id      = "${aws_subnet.public-a.id}"
  route_table_id = "${aws_route_table.public.id}"
}

resource "aws_route_table_association" "public-b" {
  subnet_id      = "${aws_subnet.public-b.id}"
  route_table_id = "${aws_route_table.public.id}"
}
