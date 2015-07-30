/*
resource "aws_route53_zone" "main" {
  name = "staging.tapglue.com"

  tags {
    Environment = "staging"
  }
}

resource "aws_route53_record" "staging" {
  zone_id = "${aws_route53_zone.main.zone_id}"
  name    = "staging.tapglue.com"
  type    = "A"

  alias {
    name                   = "${aws_elb.frontend.dns_name}"
    zone_id                = "${aws_elb.frontend.zone_id}"
    evaluate_target_health = true
  }
}
*/
