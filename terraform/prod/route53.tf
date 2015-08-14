/*
resource "aws_route53_zone" "main" {
  name = "prod.tapglue.com"

  tags {
    Environment = "prod"
  }
}

resource "aws_route53_record" "prod" {
  zone_id = "${aws_route53_zone.main.zone_id}"
  name    = "prod.tapglue.com"
  type    = "A"

  alias {
    name                   = "${aws_elb.frontend.dns_name}"
    zone_id                = "${aws_elb.frontend.zone_id}"
    evaluate_target_health = true
  }
}
*/
