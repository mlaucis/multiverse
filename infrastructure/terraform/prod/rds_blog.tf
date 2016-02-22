resource "aws_db_subnet_group" "corp-prod" {
  name        = "corp-prod"
  description = "RDS subnet group for corporate prod"
  subnet_ids  = [
    "${aws_subnet.corporate-a.id}",
    "${aws_subnet.corporate-b.id}"]
}

/**/
# Database master
resource "aws_db_instance" "corp-master" {
  identifier              = "tapglue-corp-master"
  # change this to io1 if you want to use provisioned iops for production
  storage_type            = "standard"
  #iops = 3000 # this should give us a boost in performance for production
  allocated_storage       = "5"
  engine                  = "mysql"
  engine_version          = "5.6.27"
  instance_class          = "db.t2.small"
  # if you want to change to true, see the list of instance types that support storage encryption: http://docs.aws.amazon.com/AmazonRDS/latest/UserGuide/Overview.Encryption.html#d0e10116
  storage_encrypted       = false
  name                    = "${var.rds_corp_db_name}"
  username                = "${var.rds_corp_username}"
  password                = "${var.rds_corp_password}"
  # this should be true for production
  multi_az                = true
  publicly_accessible     = false
  vpc_security_group_ids  = [
    "${aws_security_group.platform.id}",
  ]
  db_subnet_group_name    = "${aws_db_subnet_group.corp-prod.id}"
  backup_retention_period = 7
  backup_window           = "04:00-04:30"
  maintenance_window      = "sat:05:00-sat:06:30"
}
/** /
# Database slaves
resource "aws_db_instance" "corp-slave1" {
  identifier              = "corp-slave1"
  # change this to io1 if you want to use provisioned iops for production
  storage_type            = "gp2"
  #iops = 3000 # this should give us a boost in performance for production
  allocated_storage       = "5"
  engine                  = "mysql"
  engine_version          = "5.6.23"
  instance_class          = "db.t2.small"
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
  db_subnet_group_name    = "${aws_db_subnet_group.corp-prod.id}"
  backup_retention_period = 0
  backup_window           = "04:00-04:30"
  maintenance_window      = "sat:05:00-sat:06:30"
}
/**/
