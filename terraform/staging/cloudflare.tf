provider "cloudflare" {
  email = "${var.cloudflare_email}"
  token = "${var.cloudflare_token}"
}

resource "cloudflare_record" "api-staging" {
  domain = "${var.cloudflare_domain}"
  name = "api-staging"
  value = "${aws_elb.frontend.dns_name}"
  type = "CNAME"
  ttl = 1
}

resource "cloudflare_record" "dashboard-staging" {
  domain = "${var.cloudflare_domain}"
  name = "console-staging"
  value = "${aws_elb.corporate.dns_name}"
  type = "CNAME"
  ttl = 1
}

resource "cloudflare_record" "styleguide-staging" {
  domain = "${var.cloudflare_domain}"
  name = "styleguide-staging"
  value = "${aws_elb.corporate.dns_name}"
  type = "CNAME"
  ttl = 1
}

