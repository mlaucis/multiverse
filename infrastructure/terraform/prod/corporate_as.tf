resource "aws_launch_configuration" "corporate" {
  image_id                    = "${var.ami_corporate}"
  instance_type               = "t2.micro"
  associate_public_ip_address = false
  enable_monitoring           = false
  ebs_optimized               = false
  iam_instance_profile        = "${aws_iam_instance_profile.corporate.name}"

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
resource "aws_autoscaling_group" "corporate" {
  vpc_zone_identifier       = [
    "${aws_subnet.corporate-a.id}",
    "${aws_subnet.corporate-b.id}"]
  name                      = "corporate"
  max_size                  = 6
  min_size                  = 1
  health_check_type         = "EC2"
  health_check_grace_period = 60
  force_delete              = false
  launch_configuration      = "${aws_launch_configuration.corporate.name}"
  load_balancers            = [
    "${aws_elb.corporate.name}"]
  termination_policies      = [
    "OldestInstance",
    "OldestLaunchConfiguration",
    "ClosestToNextInstanceHour"]

  tag {
    key                 = "tapglue_installer"
    value               = "corporate"
    propagate_at_launch = true
  }

  tag {
    key                 = "installer_channel"
    value               = "prod"
    propagate_at_launch = true
  }

  tag {
    key                 = "Name"
    value               = "corporate"
    propagate_at_launch = true
  }
}

# Policies
resource "aws_autoscaling_policy" "corporate-increase-on-load" {
  name                   = "corporate-increase-on-load"
  scaling_adjustment     = 10
  adjustment_type        = "PercentChangeInCapacity"
  cooldown               = 60
  autoscaling_group_name = "${aws_autoscaling_group.corporate.name}"
}

resource "aws_autoscaling_policy" "corporate-decrease-on-load" {
  name                   = "corporate-decrease-on-load"
  scaling_adjustment     = 10
  adjustment_type        = "PercentChangeInCapacity"
  cooldown               = 60
  autoscaling_group_name = "${aws_autoscaling_group.corporate.name}"
}
