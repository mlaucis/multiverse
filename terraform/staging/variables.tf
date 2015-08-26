variable "vpc-region" {
  default = "eu-west-1"
}

variable "private-s3" {
  default = "com.amazonaws.eu-west-1.s3"
}

variable "ami_frontend" {
  default = "ami-0ae7bf7d"
}

variable "ami_backend" {
  default = "ami-0ae7bf7d"
}

variable "ami_corporate" {
  default = "ami-3c08534b"
}

variable "ami_bastion" {
  default = "ami-0ae7bf7d"
}

variable "ami_nat" {
  default = "ami-cb7de3bc"
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

variable "iam_profile_backend" {
  default = "staging-backend"
}

variable "iam_role_backend" {
  default = "staging-backend"
}

variable "iam_profile_frontend" {
  default = "staging-frontend"
}

variable "iam_role_frontend" {
  default = "staging-frontend"
}

variable "iam_profile_corporate" {
  default = "staging-corporate"
}

variable "iam_role_corporate" {
  default = "staging-corporate"
}
