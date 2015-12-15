resource "aws_launch_configuration" "frontend" {
  image_id                    = "${var.ami_frontend}"
  instance_type               = "t2.medium"
  associate_public_ip_address = false
  enable_monitoring           = true
  ebs_optimized               = false
  iam_instance_profile        = "${aws_iam_instance_profile.frontend.name}"

  lifecycle {
    create_before_destroy = true
  }

  security_groups             = [
    "${aws_security_group.gateway.id}",
    "${aws_security_group.private.id}",
    "${aws_security_group.service.id}",
  ]
}

# Group
resource "aws_autoscaling_group" "frontend" {
  vpc_zone_identifier       = [
    "${aws_subnet.frontend-a.id}",
    "${aws_subnet.frontend-b.id}"]
  name                      = "frontend"
  max_size                  = 30
  min_size                  = 4
  health_check_type         = "ELB"
  health_check_grace_period = 60
  force_delete              = false
  launch_configuration      = "${aws_launch_configuration.frontend.name}"
  load_balancers            = [
    "${aws_elb.frontend.name}"]
  termination_policies      = [
    "OldestInstance",
    "OldestLaunchConfiguration",
    "ClosestToNextInstanceHour"]

  tag {
    key                 = "tapglue_installer"
    value               = "intaker"
    propagate_at_launch = true
  }

  tag {
    key                 = "intaker_target"
    value               = "redis"
    propagate_at_launch = true
  }

  tag {
    key                 = "installer_channel"
    value               = "prod"
    propagate_at_launch = true
  }

  tag {
    key                 = "Name"
    value               = "frontend"
    propagate_at_launch = true
  }
}

# Policies
resource "aws_autoscaling_policy" "frontend-increase-on-load" {
  name                   = "frontend-increase-on-load"
  scaling_adjustment     = 10
  adjustment_type        = "PercentChangeInCapacity"
  cooldown               = 60
  autoscaling_group_name = "${aws_autoscaling_group.frontend.name}"
}

resource "aws_autoscaling_policy" "frontend-decrease-on-load" {
  name                   = "frontend-decrease-on-load"
  scaling_adjustment     = -10
  adjustment_type        = "PercentChangeInCapacity"
  cooldown               = 60
  autoscaling_group_name = "${aws_autoscaling_group.frontend.name}"
}

resource "aws_cloudwatch_metric_alarm" "frontend-scale-up" {
  alarm_name = "frontend-scale-up"
  comparison_operator = "GreaterThanOrEqualToThreshold"
  evaluation_periods = "2"
  metric_name = "CPUUtilization"
  namespace = "AWS/EC2"
  period = "60"
  statistic = "Average"
  threshold = "70"
  dimensions {
    AutoScalingGroupName = "${aws_autoscaling_group.frontend.name}"
  }
  alarm_description = "This metric monitor ec2 cpu utilization"
  alarm_actions = ["${aws_autoscaling_policy.frontend-increase-on-load.arn}"]
}

resource "aws_cloudwatch_metric_alarm" "frontend-scale-down" {
  alarm_name = "frontend-scale-down"
  comparison_operator = "LessThanOrEqualToThreshold"
  evaluation_periods = "2"
  metric_name = "CPUUtilization"
  namespace = "AWS/EC2"
  period = "60"
  statistic = "Average"
  threshold = "30"
  dimensions {
    AutoScalingGroupName = "${aws_autoscaling_group.frontend.name}"
  }
  alarm_description = "This metric monitor ec2 cpu utilization"
  alarm_actions = ["${aws_autoscaling_policy.frontend-decrease-on-load.arn}"]
}
