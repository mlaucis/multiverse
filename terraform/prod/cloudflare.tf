provider "cloudflare" {
  email = "${var.cloudflare_email}"
  token = "${var.cloudflare_token}"
}

resource "cloudflare_record" "api-prod" {
  domain = "${var.cloudflare_domain}"
  name = "api-prod"
  value = "${aws_elb.frontend.dns_name}"
  type = "CNAME"
  ttl = 1
}

resource "cloudflare_record" "website-prod" {
  domain = "${var.cloudflare_domain}"
  name = "website-prod"
  value = "${aws_elb.corporate.dns_name}"
  type = "CNAME"
  ttl = 1
}

resource "cloudflare_record" "dashboard-prod" {
  domain = "${var.cloudflare_domain}"
  name = "dashboard"
  value = "${aws_elb.corporate.dns_name}"
  type = "CNAME"
  ttl = 1
}

resource "cloudflare_record" "styleguide-prod" {
  domain = "${var.cloudflare_domain}"
  name = "styleguide"
  value = "${aws_elb.corporate.dns_name}"
  type = "CNAME"
  ttl = 1
}

