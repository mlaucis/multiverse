resource "aws_instance" "monitoring" {
  ami             = "${var.ami.monitoring}"
  instance_type   = "t2.medium"
  security_groups = [
    "${aws_security_group.platform.id}",
  ]
  subnet_id       = "${aws_subnet.platform-a.id}"

  tags {
    Name = "monitoring"
  }
}

resource "aws_elb" "gateway-http" {
  connection_draining         = true
  connection_draining_timeout = 10
  cross_zone_load_balancing   = true
  idle_timeout                = 30
  name                        = "gateway-http"
  security_groups             = [
    "${aws_security_group.perimeter.id}",
  ]
  subnets                     = [
    "${aws_subnet.perimeter-a.id}",
    "${aws_subnet.perimeter-b.id}",
  ]

  access_logs                 = {
    bucket    = "${aws_s3_bucket.logs-elb.id}"
    interval  = 5
  }

  health_check {
    healthy_threshold   = 2
    interval            = 5
    target              = "HTTPS:8083/health-45016490610398192"
    timeout             = 2
    unhealthy_threshold = 2
  }

  listener {
    instance_port       = 8083
    instance_protocol   = "tcp"
    lb_port             = 443
    lb_protocol         = "tcp"
  }

  tags {
    Name = "gateway-http"
  }
}

resource "aws_autoscaling_group" "service" {
  desired_capacity          = 3
  health_check_grace_period = 60
  health_check_type         = "EC2"
  launch_configuration      = "${aws_launch_configuration.service.name}"
  load_balancers            = [
    "${aws_elb.gateway-http.name}",
  ]
  max_size                  = 30
  min_size                  = 1
  name                      = "service"
  termination_policies      = [
    "OldestInstance",
    "OldestLaunchConfiguration",
    "ClosestToNextInstanceHour",
  ]
  vpc_zone_identifier       = [
    "${aws_subnet.platform-a.id}",
    "${aws_subnet.platform-b.id}",
  ]

  tag {
    key                 = "Name"
    value               = "service"
    propagate_at_launch = true
  }
}

resource "aws_launch_configuration" "service" {
  associate_public_ip_address = false
  ebs_optimized               = false
  enable_monitoring           = true
  key_name                    = "${aws_key_pair.debug.key_name}"
  iam_instance_profile        = "${aws_iam_instance_profile.ecs-agent-profile.name}"
  image_id                    = "${var.ami.ecs-agent}"
  instance_type               =  "m4.large"
  name_prefix                 = "ecs-service-"
  security_groups             = [
    "${aws_security_group.platform.id}",
  ]

  lifecycle {
    create_before_destroy = true
  }

  user_data = <<EOF
#!/bin/bash
echo ECS_CLUSTER=service >> /etc/ecs/ecs.config

# Install loggly security credentials
mkdir -pv /etc/rsyslog.d/keys/ca.d
cd /etc/rsyslog.d/keys/ca.d/
curl -O https://logdog.loggly.com/media/logs-01.loggly.com_sha12.crt

# Rsyslog for Loggly

sudo yum install -y rsyslog-gnutls
sudo mkdir -p /var/spool/rsyslog

sudo echo '$template LogglyFormat,"<%pri%>%protocol-version% %timestamp:::date-rfc3339% %HOSTNAME% %app-name% %procid% %msgid% [d2e7097f-25aa-497a-a9e3-d691bd4ec7ab@41058 tag=\"service.${var.env}.${var.region}\"] %msg%\n"' > /etc/rsyslog.d/22-loggly.conf

# Setup disk assisted queues
sudo echo '$WorkDirectory /var/spool/rsyslog # where to place spool files' >> /etc/rsyslog.d/22-loggly.conf
sudo echo '$ActionQueueFileName fwdRule1     # unique name prefix for spool files' >> /etc/rsyslog.d/22-loggly.conf
sudo echo '$ActionQueueMaxDiskSpace 100m     # 1gb space limit (use as much as possible)' >> /etc/rsyslog.d/22-loggly.conf
sudo echo '$ActionQueueSaveOnShutdown on     # save messages to disk on shutdown' >> /etc/rsyslog.d/22-loggly.conf
sudo echo '$ActionQueueType LinkedList       # run asynchronously' >> /etc/rsyslog.d/22-loggly.conf
sudo echo '$ActionResumeRetryCount -1        # infinite retries if host is down' >> /etc/rsyslog.d/22-loggly.conf

# RsyslogGnuTLS
sudo echo '$DefaultNetstreamDriverCAFile /etc/rsyslog.d/keys/ca.d/logs-01.loggly.com_sha12.crt' >> /etc/rsyslog.d/22-loggly.conf
sudo echo '$ActionSendStreamDriver gtls' >> /etc/rsyslog.d/22-loggly.conf
sudo echo '$ActionSendStreamDriverMode 1' >> /etc/rsyslog.d/22-loggly.conf
sudo echo '$ActionSendStreamDriverAuthMode x509/name' >> /etc/rsyslog.d/22-loggly.conf
sudo echo '$ActionSendStreamDriverPermittedPeer *.loggly.com' >> /etc/rsyslog.d/22-loggly.conf
sudo echo '*.* @@logs-01.loggly.com:6514;LogglyFormat' >> /etc/rsyslog.d/22-loggly.conf

sudo service rsyslog restart
EOF
}