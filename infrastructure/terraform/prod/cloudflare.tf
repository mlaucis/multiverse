resource "aws_security_group" "cloudflare-ips" {
  vpc_id      = "${aws_vpc.tapglue.id}"
  name        = "cloudflare-ips"
  description = "Cloudflare IPs"

  ingress {
    from_port   = 80
    to_port     = 80
    protocol    = "tcp"
    cidr_blocks = [
      "103.21.244.0/22",
      "103.22.200.0/22",
      "103.31.4.0/22",
      "104.16.0.0/12",
      "108.162.192.0/18",
      "141.101.64.0/18",
      "162.158.0.0/15",
      "172.64.0.0/13",
      "173.245.48.0/20",
      "188.114.96.0/20",
      "190.93.240.0/20",
      "197.234.240.0/22",
      "198.41.128.0/17",
      "199.27.128.0/21"]
  }

  ingress {
    from_port   = 443
    to_port     = 443
    protocol    = "tcp"
    cidr_blocks = [
      "103.21.244.0/22",
      "103.22.200.0/22",
      "103.31.4.0/22",
      "104.16.0.0/12",
      "108.162.192.0/18",
      "141.101.64.0/18",
      "162.158.0.0/15",
      "172.64.0.0/13",
      "173.245.48.0/20",
      "188.114.96.0/20",
      "190.93.240.0/20",
      "197.234.240.0/22",
      "198.41.128.0/17",
      "199.27.128.0/21"]
  }

  egress {
    from_port   = 80
    to_port     = 80
    protocol    = "tcp"
    cidr_blocks = [
      "103.21.244.0/22",
      "103.22.200.0/22",
      "103.31.4.0/22",
      "104.16.0.0/12",
      "108.162.192.0/18",
      "141.101.64.0/18",
      "162.158.0.0/15",
      "172.64.0.0/13",
      "173.245.48.0/20",
      "188.114.96.0/20",
      "190.93.240.0/20",
      "197.234.240.0/22",
      "198.41.128.0/17",
      "199.27.128.0/21"]
  }

  egress {
    from_port   = 443
    to_port     = 443
    protocol    = "tcp"
    cidr_blocks = [
      "103.21.244.0/22",
      "103.22.200.0/22",
      "103.31.4.0/22",
      "104.16.0.0/12",
      "108.162.192.0/18",
      "141.101.64.0/18",
      "162.158.0.0/15",
      "172.64.0.0/13",
      "173.245.48.0/20",
      "188.114.96.0/20",
      "190.93.240.0/20",
      "197.234.240.0/22",
      "198.41.128.0/17",
      "199.27.128.0/21"]
  }

  tags {
    Name = "Cloudflare IPs"
  }
}

provider "cloudflare" {
  email = "${var.cloudflare_email}"
  token = "${var.cloudflare_token}"
}

resource "cloudflare_record" "api-prod" {
  domain = "${var.cloudflare_domain}"
  name   = "api"
  value  = "${aws_elb.frontend.dns_name}"
  type   = "CNAME"
  ttl    = 1
}

resource "cloudflare_record" "website-prod" {
  domain = "${var.cloudflare_domain}"
  name   = "website-prod"
  value  = "${aws_elb.corporate.dns_name}"
  type   = "CNAME"
  ttl    = 1
}

resource "cloudflare_record" "dashboard-prod" {
  domain = "${var.cloudflare_domain}"
  name   = "dashboard"
  value  = "${aws_elb.corporate.dns_name}"
  type   = "CNAME"
  ttl    = 1
}

resource "cloudflare_record" "styleguide-prod" {
  domain = "${var.cloudflare_domain}"
  name   = "styleguide"
  value  = "${aws_elb.corporate.dns_name}"
  type   = "CNAME"
  ttl    = 1
}
