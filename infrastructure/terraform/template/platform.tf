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

resource "aws_iam_server_certificate" "monitoring" {
  name = "monitoring"
  certificate_body = "${file("${path.module}/../../certs/self/self.crt")}"
  private_key = "${file("${path.module}/../../certs/self/self.key")}"
}

resource "aws_elb" "monitoring" {
  connection_draining         = true
  connection_draining_timeout = 10
  cross_zone_load_balancing   = true
  idle_timeout                = 30
  name                        = "monitoring"
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

  instances = [
    "${aws_instance.monitoring.id}",
  ]

  listener {
    instance_port       = 3000
    instance_protocol   = "http"
    lb_port             = 443
    lb_protocol         = "https"
    ssl_certificate_id  = "${aws_iam_server_certificate.monitoring.arn}"
  }

  tags {
    Name = "monitoring"
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

  user_data                   = <<EOF
#!/bin/bash
echo ECS_CLUSTER=service >> /etc/ecs/ecs.config

# Install loggly security credentials
mkdir -pv /etc/rsyslog.d/keys/ca.d
cd /etc/rsyslog.d/keys/ca.d/
curl -O https://logdog.loggly.com/media/logs-01.loggly.com_sha12.crt

# Rsyslog for Loggly

sudo yum install -y rsyslog-gnutls
sudo mkdir -p /var/spool/rsyslog

echo '$template LogglyFormat,"<%pri%>%protocol-version% %timestamp:::date-rfc3339% %HOSTNAME% %app-name% %procid% %msgid% [d2e7097f-25aa-497a-a9e3-d691bd4ec7ab@41058 tag=\"service.prod.eu-central-1\"] %msg%\n"

# Setup disk assisted queues
$WorkDirectory /var/spool/rsyslog # where to place spool files
$ActionQueueFileName fwdRule1     # unique name prefix for spool files
$ActionQueueMaxDiskSpace 100m     # 1gb space limit (use as much as possible)
$ActionQueueSaveOnShutdown on     # save messages to disk on shutdown
$ActionQueueType LinkedList       # run asynchronously
$ActionResumeRetryCount -1        # infinite retries if host is down

# RsyslogGnuTLS
$DefaultNetstreamDriverCAFile /etc/rsyslog.d/keys/ca.d/logs-01.loggly.com_sha12.crt
$ActionSendStreamDriver gtls
$ActionSendStreamDriverMode 1
$ActionSendStreamDriverAuthMode x509/name
$ActionSendStreamDriverPermittedPeer *.loggly.com
*.* @@logs-01.loggly.com:6514;LogglyFormat
' | sudo tee /etc/rsyslog/22-loggly.conf > /dev/null

sudo service rsyslog restart

echo '#!/bin/sh

/usr/sbin/logrotate /etc/logrotate.hourly.conf >/dev/null 2>&1
EXITVALUE=$?
if [ $EXITVALUE != 0 ]; then
    /usr/bin/logger -t logrotate "ALERT exited abnormally with [$EXITVALUE]"
fi
exit 0
' | sudo tee /etc/cron.hourly/logrotate > /dev/null

sudo chmod +x /etc/cron.hourly/logrotate

echo '/var/log/messages {
    compress
    create
    daily
    rotate 5
    size 100M
    postrotate
  /bin/kill -HUP `cat /var/run/syslogd.pid 2> /dev/null` 2> /dev/null || true
    endscript
}' | sudo tee /etc/logrotate.hourly.conf > /dev/null

EOF
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