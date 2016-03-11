resource "cloudflare_record" "adrenalynmapp" {
  domain  = "${var.cloudflare_domain}"
  name    = "adrenalynmapp"
  value   = "${aws_elb.frontend.dns_name}"
  type    = "CNAME"
  ttl     = 1
  proxied = true
}

resource "cloudflare_record" "dailyme" {
  domain  = "${var.cloudflare_domain}"
  name    = "dailyme"
  value   = "${aws_elb.gateway-http.dns_name}"
  type    = "CNAME"
  ttl     = 1
  proxied = true
}

resource "cloudflare_record" "dawanda" {
  domain  = "${var.cloudflare_domain}"
  name    = "dawanda"
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
