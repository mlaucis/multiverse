resource "aws_route53_zone" "tapglue-internal" {
  name    = "tapglue.int"
  vpc_id  = "${aws_vpc.tapglue.id}"
  comment = "Internal prod zone"

  tags {
    Environment = "prod"
  }
}

resource "aws_route53_zone" "env" {
  comment = "zone to isolate DNS routes for ${var.env}.${var.region}"
  name    = "${var.env}.${var.region}"
  vpc_id  = "${aws_vpc.tapglue.id}"
}

resource "aws_route53_record" "ratelimiter-cache" {
  name    = "cache.ratelimiter"
  ttl     = "5"
  type    = "CNAME"
  zone_id = "${aws_route53_zone.env.id}"
  records = [
    "${aws_elasticache_cluster.rate-limiter.cache_nodes.0.address}",
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

resource "aws_route53_record" "db-master" {
  zone_id = "${aws_route53_zone.tapglue-internal.zone_id}"
  name    = "db-master"
  type    = "CNAME"

  ttl     = "5"
  records = [
    "${aws_db_instance.master.address}",
  ]
}

resource "aws_route53_record" "db-slave1" {
  zone_id = "${aws_route53_zone.tapglue-internal.zone_id}"
  name    = "db-slave1"
  type    = "CNAME"

  ttl     = "5"
  records = [
    "${aws_db_instance.master.address}"]
}

resource "aws_route53_record" "db-slave2" {
  zone_id = "${aws_route53_zone.tapglue-internal.zone_id}"
  name    = "db-slave2"
  type    = "CNAME"

  ttl     = "5"
  records = [
    "${aws_db_instance.master.address}"]
}


/*
resource "aws_route53_record" "db-corp-master" {
  zone_id = "${aws_route53_zone.tapglue-internal.zone_id}"
  name    = "db-corp-master"
  type    = "CNAME"

  ttl     = "5"
  records = [
    "${aws_db_instance.corp-master.address}"]
}

resource "aws_route53_record" "db-corp-slave1" {
  zone_id = "${aws_route53_zone.tapglue-internal.zone_id}"
  name    = "db-corp-slave1"
  type    = "CNAME"

  ttl     = "5"
  records = [
    "${aws_db_instance.corp-master.address}"]
}
*/

resource "aws_route53_record" "rate-limiter" {
  zone_id = "${aws_route53_zone.tapglue-internal.zone_id}"
  name    = "rate-limiter"
  type    = "CNAME"

  ttl     = "5"
  records = [
    "${aws_elasticache_cluster.rate-limiter.cache_nodes.0.address}",
  ]
}

resource "aws_route53_record" "cache-app" {
  zone_id = "${aws_route53_zone.tapglue-internal.zone_id}"
  name    = "cache-app"
  type    = "CNAME"

  ttl     = "5"
  records = [
    "${aws_elasticache_cluster.rate-limiter.cache_nodes.0.address}"]
}
