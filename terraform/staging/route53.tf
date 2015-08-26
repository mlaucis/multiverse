/** /
resource "aws_route53_delegation_set" "staging" {
  reference_name = "custom"
}

resource "aws_route53_zone" "staging" {
  name = "staging.tapglue.com"
  delegation_set_id = "${aws_route53_delegation_set.staging.id}"

  tags {
    Environment = "staging"
  }
}

resource "aws_route53_record" "api-staging" {
  zone_id = "${aws_route53_zone.staging.zone_id}"
  name    = "api.staging.tapglue.com"
  type    = "A"

  alias {
    name                   = "${aws_elb.frontend.dns_name}"
    zone_id                = "${aws_elb.frontend.zone_id}"
    evaluate_target_health = true
  }
}

resource "aws_route53_record" "dashboard-staging" {
  zone_id = "${aws_route53_zone.staging.zone_id}"
  name    = "dashboard.staging.tapglue.com"
  type    = "A"

  alias {
    name                   = "${aws_elb.corporate.dns_name}"
    zone_id                = "${aws_elb.corporate.zone_id}"
    evaluate_target_health = true
  }
}
/**/
