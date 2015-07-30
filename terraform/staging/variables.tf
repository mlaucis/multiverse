variable "vpc-region" {
  default = "eu-west-1"
}

variable "private-s3" {
  default = "com.amazonaws.eu-west-1.s3"
}

variable "zone-a" {
  default = "eu-west-1a"
}

variable "zone-b" {
  default = "eu-west-1b"
}

variable "zone-bastion" {
  default = "eu-west-1a"
}

variable "zone-nat" {
  default = "eu-west-1a"
}

variable "bastion-size" {
  default = "t2.micro"
}

variable "nat-size" {
  default = "m1.small"
}

variable "aws_key_path" {
  default = ""
}

variable "aws_key_name" {
  default = "demoterra"
}

variable "aws_nat_ami" {
  default = "ami-cb7de3bc"
}

variable "aws_ubuntu_ami" {
  default = "ami-7e6b3f09"
}
