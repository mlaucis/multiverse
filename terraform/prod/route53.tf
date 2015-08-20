/** /
resource "aws_route53_delegation_set" "prod" {
  reference_name = "custom"
}

resource "aws_route53_zone" "prod" {
  name = "prod.tapglue.com"
  delegation_set_id = "${aws_route53_delegation_set.prod.id}"

  tags {
    Environment = "prod"
  }
}

resource "aws_route53_record" "api-prod" {
  zone_id = "${aws_route53_zone.prod.zone_id}"
  name    = "api.prod.tapglue.com"
  type    = "A"

  alias {
    name                   = "${aws_elb.frontend.dns_name}"
    zone_id                = "${aws_elb.frontend.zone_id}"
    evaluate_target_health = true
  }
}

resource "aws_route53_record" "dashboard-prod" {
  zone_id = "${aws_route53_zone.prod.zone_id}"
  name    = "dashboard.prod.tapglue.com"
  type    = "A"

  alias {
    name                   = "${aws_elb.corporate.dns_name}"
    zone_id                = "${aws_elb.corporate.zone_id}"
    evaluate_target_health = true
  }
}
/**/
