resource "aws_route53_zone" "tapglue-internal" {
  name    = "tapglue.int"
  vpc_id  = "${aws_vpc.tapglue.id}"
  comment = "Internal prod zone"

  tags {
    Environment = "prod"
  }
}

resource "aws_route53_record" "db-master" {
  zone_id = "${aws_route53_zone.tapglue-internal.zone_id}"
  name    = "db-master"
  type    = "CNAME"

  ttl     = "5"
  records = [
    "${aws_db_instance.master.address}"]
}

resource "aws_route53_record" "db-slave1" {
  zone_id = "${aws_route53_zone.tapglue-internal.zone_id}"
  name    = "db-slave1"
  type    = "CNAME"

  ttl     = "5"
  records = [
    "${aws_db_instance.master.address}"]
}

resource "aws_route53_record" "rate-limiter" {
  zone_id = "${aws_route53_zone.tapglue-internal.zone_id}"
  name    = "rate-limiter"
  type    = "CNAME"

  ttl     = "5"
  records = [
    "${aws_elasticache_cluster.rate-limiter.cache_nodes.0.address}"]
}

resource "aws_route53_record" "cache-app" {
  zone_id = "${aws_route53_zone.tapglue-internal.zone_id}"
  name    = "cache-app"
  type    = "CNAME"

  ttl     = "5"
  records = [
    "${aws_elasticache_cluster.rate-limiter.cache_nodes.0.address}"]
}
