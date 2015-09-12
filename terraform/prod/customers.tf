resource "cloudflare_record" "dailyme" {
  domain = "${var.cloudflare_domain}"
  name   = "dailyme"
  value  = "${aws_elb.frontend.dns_name}"
  type   = "CNAME"
  ttl    = 1
}

resource "cloudflare_record" "dawanda" {
  domain = "${var.cloudflare_domain}"
  name   = "dawanda"
  value  = "${aws_elb.frontend.dns_name}"
  type   = "CNAME"
  ttl    = 1
}
