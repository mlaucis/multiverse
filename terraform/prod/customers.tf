resource "cloudflare_record" "dailyme" {
  domain = "${var.cloudflare_domain}"
  name = "dailyme"
#  value = "${aws_elb.frontend.dns_name}"
  value = "staging-347078730.eu-central-1.elb.amazonaws.com"
  type = "CNAME"
  ttl = 1
}

resource "cloudflare_record" "dawanda" {
  domain = "${var.cloudflare_domain}"
  name = "dawanda"
  #  value = "${aws_elb.frontend.dns_name}"
  value = "staging-347078730.eu-central-1.elb.amazonaws.com"
  type = "CNAME"
  ttl = 1
}
