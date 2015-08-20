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

resource "cloudflare_record" "dashboard-prod" {
  domain = "${var.cloudflare_domain}"
  name = "dashboard-prod"
  value = "${aws_elb.corporate.dns_name}"
  type = "CNAME"
  ttl = 1
}
