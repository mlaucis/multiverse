provider "cloudflare" {
  email = "${var.cloudflare_email}"
  token = "${var.cloudflare_token}"
}

resource "cloudflare_record" "api-prod" {
  domain  = "${var.cloudflare_domain}"
  name    = "api"
  value   = "${aws_elb.gateway-http.dns_name}"
  type    = "CNAME"
  ttl     = 1
  proxied = true
}

resource "cloudflare_record" "website-prod" {
  domain  = "${var.cloudflare_domain}"
  name    = "website-prod"
  value   = "tapglue.github.io"
  type    = "CNAME"
  ttl     = 1
  proxied = true
}

resource "cloudflare_record" "monitoring" {
  domain  = "${var.cloudflare_domain}"
  name    = "monitoring-${var.env}-${var.region}"
  proxied = true
  ttl     = 1
  type    = "CNAME"
  value   = "${aws_elb.monitoring.dns_name}"
}

# FIX: Temporary public endpoint for monitoring dashboard.
resource "cloudflare_record" "monitoring-old" {
  domain  = "${var.cloudflare_domain}"
  name    = "monitoring"
  ttl     = 1
  type    = "CNAME"
  value   = "${aws_elb.monitoring.dns_name}"
  proxied = true
}
