resource "aws_subnet" "rds-a" {
  vpc_id                  = "${aws_vpc.tapglue.id}"
  map_public_ip_on_launch = false

  cidr_block              = "10.0.36.0/22"
  availability_zone       = "${var.zone-a}"

  tags {
    Name = "RDS A"
  }
}

resource "aws_subnet" "rds-b" {
  vpc_id                  = "${aws_vpc.tapglue.id}"
  map_public_ip_on_launch = false

  cidr_block              = "10.0.40.0/22"
  availability_zone       = "${var.zone-b}"

  tags {
    Name = "RDS B"
  }
}

resource "aws_db_subnet_group" "prod" {
  name        = "prod"
  description = "RDS subnet group for prod"
  subnet_ids  = [
    "${aws_subnet.rds-a.id}",
    "${aws_subnet.rds-b.id}"]
}

resource "aws_db_parameter_group" "master-prod" {
  name        = "master-prod"
  family      = "postgres9.4"
  description = "Postgres prod parameter group"

  parameter {
    name         = "log_statement"
    value        = "all"
    apply_method = "immediate"
  }

  parameter {
    name         = "log_min_duration_statement"
    value        = "20"
    apply_method = "immediate"
  }

  parameter {
    name         = "log_duration"
    value        = "1"
    apply_method = "immediate"
  }

  parameter {
    name         = "shared_preload_libraries"
    value        = "pg_stat_statements"
    apply_method = "immediate"
  }

  parameter {
    name         = "track_activity_query_size"
    value        = "2048"
    apply_method = "immediate"
  }

  parameter {
    name         = "pg_stat_statements.track"
    value        = "ALL"
    apply_method = "immediate"
  }

  parameter {
    name         = "autovacuum"
    value        = "1"
    apply_method = "immediate"
  }

  parameter {
    name         = "log_autovacuum_min_duration"
    value        = "1"
    apply_method = "immediate"
  }

  parameter {
    name         = "checkpoint_segments"
    value        = "8"
    apply_method = "immediate"
  }

  parameter {
    name         = "max_connections"
    value        = "128"
    apply_method = "immediate"
  }

  parameter {
    name         = "maintenance_work_mem"
    value        = "768000"
    apply_method = "immediate"
  }

  parameter {
    name         = "work_mem"
    value        = "64000"
    apply_method = "immediate"
  }
}

# Database master
resource "aws_db_instance" "master" {
  identifier              = "${var.rds_id}"
  # change this to io1 if you want to use provisioned iops for production
  storage_type            = "io1"
  iops                    = 1000 # this should give us a boost in performance for production
  allocated_storage       = "200"
  engine                  = "postgres"
  engine_version          = "9.4.7"
  instance_class          = "db.r3.large"
  # if you want to change to true, see the list of instance types that support storage encryption: http://docs.aws.amazon.com/AmazonRDS/latest/UserGuide/Overview.Encryption.html#d0e10116
  storage_encrypted       = true
  name                    = "${var.rds_db_name}"
  username                = "${var.rds_username}"
  password                = "${var.rds_password}"
  multi_az                = true
  monitoring_role_arn     = "${var.rds_monitoring_role}"
  monitoring_interval     = 1
  # this should be true for production
  publicly_accessible     = false
  vpc_security_group_ids  = [
    "${aws_security_group.platform.id}",
  ]
  db_subnet_group_name    = "${aws_db_subnet_group.prod.id}"
  backup_retention_period = 30
  backup_window           = "00:00-01:30"
  maintenance_window      = "tue:02:00-tue:03:00"
  parameter_group_name    = "${aws_db_parameter_group.master-prod.id}"
  apply_immediately       = true
  skip_final_snapshot     = false
}

# Database slaves
/** /
resource "aws_db_instance" "slave1" {
  identifier              = "slave1"
  # change this to io1 if you want to use provisioned iops for production
  storage_type            = "gp2"
  #iops = 3000 # this should give us a boost in performance for production
  allocated_storage       = "200"
  engine                  = "postgres"
  engine_version          = "9.4.4"
  instance_class          = "db.r3.xlarge"
  # if you want to change to true, see the list of instance types that support storage encryption: http://docs.aws.amazon.com/AmazonRDS/latest/UserGuide/Overview.Encryption.html#d0e10116
  storage_encrypted       = true
  name                    = "${var.rds_db_name}"
  username                = "${var.rds_username}"
  password                = "${var.rds_password}"
  multi_az                = false
  publicly_accessible     = false
  replicate_source_db     = "${aws_db_instance.master.identifier}"
  vpc_security_group_ids  = [
    "${aws_security_group.platform.id}",
  ]
  db_subnet_group_name    = "${aws_db_subnet_group.prod.id}"
  backup_retention_period = 0
  backup_window           = "04:00-04:30"
  maintenance_window      = "sat:05:00-sat:06:30"
  parameter_group_name    = "${aws_db_parameter_group.master-prod.id}"
  apply_immediately       = true
}
/** /
resource "aws_db_instance" "slave2" {
  identifier              = "slave2"
  # change this to io1 if you want to use provisioned iops for production
  storage_type            = "gp2"
  #iops = 3000 # this should give us a boost in performance for production
  allocated_storage       = "100"
  engine                  = "postgres"
  engine_version          = "9.4.4"
  instance_class          = "db.r3.xlarge"
  # if you want to change to true, see the list of instance types that support storage encryption: http://docs.aws.amazon.com/AmazonRDS/latest/UserGuide/Overview.Encryption.html#d0e10116
  storage_encrypted       = true
  name                    = "${var.rds_db_name}"
  username                = "${var.rds_username}"
  password                = "${var.rds_password}"
  multi_az                = false
  publicly_accessible     = false
  replicate_source_db     = "${aws_db_instance.master.identifier}"
  vpc_security_group_ids  = [
    "${aws_security_group.platform.id}",
  ]
  db_subnet_group_name    = "${aws_db_subnet_group.prod.id}"
  backup_retention_period = 0
  backup_window           = "04:00-04:30"
  maintenance_window      = "sat:05:00-sat:06:30"
  parameter_group_name    = "${aws_db_parameter_group.master-prod.id}"
  apply_immediately       = true
}
/**/
