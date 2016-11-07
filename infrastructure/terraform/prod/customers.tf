resource "cloudflare_record" "avakin" {
  domain = "${var.cloudflare_domain}"
  name    = "avakin"
  value   = "${aws_elb.gateway-http.dns_name}"
  type    = "CNAME"
  ttl     = 1
  proxied = true
}

resource "cloudflare_record" "bikestorming" {
  domain = "${var.cloudflare_domain}"
  name    = "bikestorming"
  value   = "${aws_elb.gateway-http.dns_name}"
  type    = "CNAME"
  ttl     = 1
  proxied = true
}

resource "cloudflare_record" "gambify" {
  domain  = "${var.cloudflare_domain}"
  name    = "gambify"
  value   = "${aws_elb.gateway-http.dns_name}"
  type    = "CNAME"
  ttl     = 1
  proxied = true
}

resource "cloudflare_record" "stepz" {
  domain  = "${var.cloudflare_domain}"
  name    = "stepz"
  value   = "${aws_elb.gateway-http.dns_name}"
  type    = "CNAME"
  ttl     = 1
  proxied = true
}

resource "cloudflare_record" "umake" {
  domain  = "${var.cloudflare_domain}"
  name    = "umake"
  value   = "${aws_elb.gateway-http.dns_name}"
  type    = "CNAME"
  ttl     = 1
  proxied = true
}
