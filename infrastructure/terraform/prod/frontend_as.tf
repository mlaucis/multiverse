resource "aws_launch_configuration" "frontend" {
  depends_on                  = [
    "aws_security_group.frontend-elb-vpc",
    "aws_security_group.frontend-ssh"]

  image_id                    = "${var.ami_frontend}"
  instance_type               = "c3.large"
  associate_public_ip_address = false
  enable_monitoring           = true
  ebs_optimized               = false
  iam_instance_profile        = "${aws_iam_instance_profile.frontend.name}"

  lifecycle {
    create_before_destroy = true
  }

  security_groups             = [
    "${aws_security_group.frontend-elb-vpc.id}",
    "${aws_security_group.frontend-ssh.id}",
    "${aws_security_group.to-nat.id}",
    "${aws_security_group.rds_ec2.id}",
    "${aws_security_group.ec-redis-ec2.id}",
  ]
}

# Group
resource "aws_autoscaling_group" "frontend" {
  vpc_zone_identifier       = [
    "${aws_subnet.frontend-a.id}",
    "${aws_subnet.frontend-b.id}"]
  name                      = "frontend"
  max_size                  = 30
  min_size                  = 2
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
    value               = "kinesis"
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
  scaling_adjustment     = 10
  adjustment_type        = "PercentChangeInCapacity"
  cooldown               = 60
  autoscaling_group_name = "${aws_autoscaling_group.frontend.name}"
}
