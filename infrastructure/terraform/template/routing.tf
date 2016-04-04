variable "domain" {
  default     = "tapglue.com"
  description = "Publicly addressable domain"
  type        = "string"
}

resource "aws_route53_zone" "env" {
  comment = "zone to isolate DNS routes for ${var.env}.${var.region}"
  name    = "${var.env}.${var.region}"
  vpc_id  = "${aws_vpc.env.id}"
}

resource "aws_route53_record" "ratelimiter-cache" {
  name    = "cache.ratelimiter"
  ttl     = "5"
  type    = "CNAME"
  zone_id = "${aws_route53_zone.env.id}"
  records = [
    "${aws_elasticache_cluster.ratelimiter.cache_nodes.0.address}",
  ]
}

resource "aws_route53_record" "service-db" {
  name    = "db-master.service"
  ttl     = "5"
  type    = "CNAME"
  zone_id = "${aws_route53_zone.env.id}"
  records = [
    "${aws_db_instance.master.address}",
  ]
}

provider "cloudflare" {
  email = "${var.cloudflare.email}"
  token = "${var.cloudflare.token}"
}

resource "cloudflare_record" "api" {
  domain  = "${var.domain}"
  name    = "api-${var.env}-${var.region}"
  proxied = true
  ttl     = 1
  type    = "CNAME"
  value   = "${aws_elb.gateway-http.dns_name}"
}

resource "cloudflare_record" "monitoring" {
  domain  = "${var.domain}"
  name    = "monitoring-${var.env}-${var.region}"
  proxied = true
  ttl     = 1
  type    = "CNAME"
  value   = "${aws_elb.monitoring.dns_name}"
}