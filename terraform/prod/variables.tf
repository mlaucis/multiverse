variable "vpc-region" {
  default = "eu-central-1"
}

variable "private-s3" {
  default = "com.amazonaws.eu-central-1.s3"
}

variable "ami_frontend" {
  default = "ami-94e5e389"
}

variable "ami_backend" {
  default = "ami-94e5e389"
}

variable "ami_corporate" {
  default = "ami-aa3031b7"
}

variable "ami_bastion" {
  default = "ami-94e5e389"
}

variable "ami_nat" {
  default = "ami-1e073a03"
}

variable "zone-a" {
  default = "eu-central-1a"
}

variable "zone-b" {
  default = "eu-central-1b"
}

variable "zone-bastion" {
  default = "eu-central-1a"
}

variable "zone-nat" {
  default = "eu-central-1a"
}

variable "bastion-size" {
  default = "t2.micro"
}

variable "nat-size" {
  default = "t2.small"
}

variable "iam_profile_backend" {
  default = "prod-backend"
}

variable "iam_role_backend" {
  default = "prod-backend"
}

variable "iam_profile_frontend" {
  default = "prod-frontend"
}

variable "iam_role_frontend" {
  default = "prod-frontend"
}

variable "iam_profile_corporate" {
  default = "prod-corporate"
}

variable "iam_role_corporate" {
  default = "prod-corporate"
}

variable "cloudflare_email" {
  default = "tools@tapglue.com"
}

variable "cloudflare_token" {
  default = "8495c1d8eadae7413a79f74fa3bd3116c8c1b"
}

variable "cloudflare_domain" {
  default = "tapglue.com"
}

variable "rds_db_name" {
  default = "tapglue"
}

variable "rds_username" {
  default = "tapglue"
}

variable "rds_password" {
  default = "gFthJy858iXOIA3IM0GARIuFYdIWkeCHJc0vto"
}
