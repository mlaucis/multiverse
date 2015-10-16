resource "aws_launch_configuration" "backend" {
  depends_on                  = [
    "aws_security_group.backend-ssh"]

  image_id                    = "${var.ami_backend}"
  instance_type               = "c4.large"
  associate_public_ip_address = false
  enable_monitoring           = true
  ebs_optimized               = false
  iam_instance_profile        = "${aws_iam_instance_profile.backend.name}"

  lifecycle {
    create_before_destroy = true
  }

  security_groups             = [
    "${aws_security_group.backend-ssh.id}",
    "${aws_security_group.to-nat.id}",
    "${aws_security_group.rds_ec2.id}",
    "${aws_security_group.ec-redis-ec2.id}",
  ]
}

# Group
resource "aws_autoscaling_group" "backend" {
  vpc_zone_identifier       = [
    "${aws_subnet.backend-a.id}",
    "${aws_subnet.backend-b.id}"]
  name                      = "backend"
  max_size                  = 1
  min_size                  = 1
  health_check_type         = "EC2"
  health_check_grace_period = 60
  force_delete              = false
  launch_configuration      = "${aws_launch_configuration.backend.name}"
  termination_policies      = [
    "OldestInstance",
    "OldestLaunchConfiguration",
    "ClosestToNextInstanceHour"]

  tag {
    key                 = "tapglue_installer"
    value               = "distributor"
    propagate_at_launch = true
  }

  tag {
    key                 = "distributor_target"
    value               = "postgres"
    propagate_at_launch = true
  }

  tag {
    key                 = "installer_channel"
    value               = "prod"
    propagate_at_launch = true
  }

  tag {
    key                 = "Name"
    value               = "backend"
    propagate_at_launch = true
  }
}

# Policies
resource "aws_autoscaling_policy" "backend-increase-on-load" {
  name                   = "backend-increase-on-load"
  scaling_adjustment     = 10
  adjustment_type        = "PercentChangeInCapacity"
  cooldown               = 60
  autoscaling_group_name = "${aws_autoscaling_group.backend.name}"
}

resource "aws_autoscaling_policy" "backend-decrease-on-load" {
  name                   = "backend-decrease-on-load"
  scaling_adjustment     = 10
  adjustment_type        = "PercentChangeInCapacity"
  cooldown               = 60
  autoscaling_group_name = "${aws_autoscaling_group.backend.name}"
}
