resource "aws_autoscaling_notification" "autoscaling-staging" {
  group_names   = [
    "${aws_autoscaling_group.frontend.name}",
    "${aws_autoscaling_group.backend.name}",
  ]
  notifications = [
    "autoscaling:EC2_INSTANCE_LAUNCH",
    "autoscaling:EC2_INSTANCE_LAUNCH_ERROR",
    "autoscaling:EC2_INSTANCE_TERMINATE",
    "autoscaling:EC2_INSTANCE_TERMINATE_ERROR"]

  topic_arn     = "${aws_sns_topic.autoscaling-staging.arn}"
}

resource "aws_sns_topic" "autoscaling-staging" {
  name = "autoscaling-staging"
}

# TODO Since Terraform doesn't support email as subcription so we have to do it manually
/*resource "aws_sns_topic_subscription" "user_updates_sqs_target" {
  topic_arn = "${aws_sns_topic.autoscaling.arn}"
  protocol  = "email"
  endpoint  = "alerts@tapglue.com"
}*/
