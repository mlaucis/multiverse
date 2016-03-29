resource "aws_db_subnet_group" "master" {
  description = "Postgres master"
  name = "master"
  subnet_ids  = [
    "${aws_subnet.platform-a.id}",
    "${aws_subnet.platform-b.id}",
  ]
}

resource "aws_db_parameter_group" "master" {
  description = "Postgres master"
  family      = "postgres9.4"
  name        = "master"

  parameter {
    apply_method  = "pending-reboot"
    name          = "log_statement"
    value         = "all"
  }

  parameter {
    apply_method  = "pending-reboot"
    name          = "log_min_duration_statement"
    value         = "20"
  }

  parameter {
    apply_method  = "pending-reboot"
    name          = "log_duration"
    value         = "1"
  }

  parameter {
    apply_method  = "pending-reboot"
    name          = "shared_preload_libraries"
    value         = "pg_stat_statements"
  }

  parameter {
    apply_method  = "pending-reboot"
    name          = "track_activity_query_size"
    value         = "2048"
  }

  parameter {
    apply_method  = "pending-reboot"
    name          = "pg_stat_statements.track"
    value         = "ALL"
  }

  parameter {
    apply_method  = "pending-reboot"
    name          = "autovacuum"
    value         = "1"
  }

  parameter {
    apply_method  = "pending-reboot"
    name          = "log_autovacuum_min_duration"
    value         = "1"
  }
}

resource "aws_db_instance" "master" {
  allocated_storage         = "200"
  apply_immediately         = true
  backup_retention_period   = 30
  backup_window             = "04:00-04:30"
  db_subnet_group_name      = "${aws_db_subnet_group.master.id}"
  final_snapshot_identifier = "db-master-${var.env}-${var.region}-final"
  identifier                = "master"
  iops                      = 1000
  storage_type              = "io1"
  engine                    = "postgres"
  engine_version            = "9.4.5"
  instance_class            = "db.r3.large"
  maintenance_window        = "sat:05:00-sat:06:30"
  monitoring_interval       = 1
  monitoring_role_arn       = "${var.role.rds-monitoring-role}"
  multi_az                  = true
  parameter_group_name      = "${aws_db_parameter_group.master.id}"
  publicly_accessible       = false
  skip_final_snapshot       = false
  storage_encrypted         = true
  vpc_security_group_ids    = [
    "${aws_security_group.platform.id}",
  ]

  name                    = "${var.pg_db_name}"
  username                = "${var.pg_username}"
  password                = "${var.pg_password}"
}

resource "aws_elasticache_subnet_group" "ratelimiter" {
  description = "ratelimiter cache"
  name        = "ratelimiter"
  subnet_ids  = [
    "${aws_subnet.platform-a.id}",
    "${aws_subnet.platform-b.id}",
  ]
}

resource "aws_elasticache_cluster" "ratelimiter" {
  cluster_id            = "ratelimiter"
  engine                = "redis"
  engine_version        = "2.8.21"
  maintenance_window    = "sun:05:00-sun:06:00"
  node_type             = "cache.t2.micro"
  num_cache_nodes       = 1
  parameter_group_name  = "default.redis2.8"
  port                  = 6379
  security_group_ids    = [
    "${aws_security_group.platform.id}",
  ]
  subnet_group_name     = "${aws_elasticache_subnet_group.ratelimiter.name}"
}

resource "aws_s3_bucket" "logs-elb" {
  bucket        = "${var.region}-${var.env}-logs-elb"
  force_destroy = true
  policy        = <<EOF
{
	"Version": "2012-10-17",
	"Id": "Policy1458936351610",
	"Statement": [
		{
			"Sid": "Stmt1458936348932",
			"Effect": "Allow",
			"Principal": {
				"AWS": "arn:aws:iam::127311923021:root"
			},
			"Action": "s3:PutObject",
			"Resource": "arn:aws:s3:::us-east-1-test-logs-elb/AWSLogs/775034650473/*"
		}
	]
}
EOF
}