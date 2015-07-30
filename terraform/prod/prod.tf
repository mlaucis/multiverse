provider "aws" {
  region = "${var.vpc-region}"
}

resource "aws_vpc" "prod" {
  cidr_block           = "10.0.0.0/16"
  enable_dns_support   = true
  enable_dns_hostnames = true

  tags {
    Name = "Production - Tapglue"
  }
}

resource "aws_internet_gateway" "default" {
  vpc_id = "${aws_vpc.prod.id}"
}

resource "aws_vpc_endpoint" "private-s3" {
  vpc_id          = "${aws_vpc.prod.id}"
  service_name    = "${var.private-s3}"
  route_table_ids = [
    "${aws_route_table.to-nat.id}"]
}
