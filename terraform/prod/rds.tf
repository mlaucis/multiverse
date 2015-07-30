resource "aws_subnet" "rds-a" {
  vpc_id                  = "${aws_vpc.prod.id}"
  map_public_ip_on_launch = false

  cidr_block              = "10.0.36.0/22"
  availability_zone       = "${var.zone-a}"

  tags {
    Name = "RDS A"
  }
}

resource "aws_subnet" "rds-b" {
  vpc_id                  = "${aws_vpc.prod.id}"
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

# Security groups
resource "aws_security_group" "rds_db" {
  depends_on  = [
    "aws_db_subnet_group.prod"]
  name        = "RDS incoming traffic"
  description = "Allow traffic on postgres port only"
  vpc_id      = "${aws_vpc.prod.id}"

  ingress {
    from_port   = 5432
    to_port     = 5432
    protocol    = "tcp"
    cidr_blocks = [
      "${aws_subnet.frontend-a.cidr_block}"]
  }
  ingress {
    from_port   = 5432
    to_port     = 5432
    protocol    = "tcp"
    cidr_blocks = [
      "${aws_subnet.frontend-b.cidr_block}"]
  }

  ingress {
    from_port   = 5432
    to_port     = 5432
    protocol    = "tcp"
    cidr_blocks = [
      "${aws_subnet.backend-a.cidr_block}"]
  }
  ingress {
    from_port   = 5432
    to_port     = 5432
    protocol    = "tcp"
    cidr_blocks = [
      "${aws_subnet.backend-b.cidr_block}"]
  }

  egress {
    from_port   = 5432
    to_port     = 5432
    protocol    = "tcp"
    cidr_blocks = [
      "${aws_subnet.frontend-a.cidr_block}"]
  }
  egress {
    from_port   = 5432
    to_port     = 5432
    protocol    = "tcp"
    cidr_blocks = [
      "${aws_subnet.frontend-b.cidr_block}"]
  }

  egress {
    from_port   = 5432
    to_port     = 5432
    protocol    = "tcp"
    cidr_blocks = [
      "${aws_subnet.backend-a.cidr_block}"]
  }
  egress {
    from_port   = 5432
    to_port     = 5432
    protocol    = "tcp"
    cidr_blocks = [
      "${aws_subnet.backend-b.cidr_block}"]
  }

  tags {
    Name = "RDS incoming traffic"
  }
}

resource "aws_security_group" "rds_ec2" {
  depends_on  = [
    "aws_db_subnet_group.prod"]
  name        = "RDS outgoing traffic"
  description = "Allow traffic to postgres port only"
  vpc_id      = "${aws_vpc.prod.id}"

  egress {
    from_port       = 5432
    to_port         = 5432
    protocol        = "tcp"
    security_groups = [
      "${aws_security_group.rds_db.id}"]
  }
  ingress {
    from_port       = 5432
    to_port         = 5432
    protocol        = "tcp"
    security_groups = [
      "${aws_security_group.rds_db.id}"]
  }

  tags {
    Name = "RDS outgoing traffic"
  }
}

/**/
# Database master
resource "aws_db_instance" "master" {
  identifier              = "tapglue-master"
  storage_type            = "standard"
  # change this to io1 if you want to use provisioned iops for production
  #iops = 3000 # this should give us a boost in performance for production
  allocated_storage       = "10"
  engine                  = "postgres"
  engine_version          = "9.4.4"
  instance_class          = "db.t2.micro"
  storage_encrypted       = false
  # if you want to change to true, see the list of instance types that support storage encryption: http://docs.aws.amazon.com/AmazonRDS/latest/UserGuide/Overview.Encryption.html#d0e10116
  name                    = "tapglue_prod"
  username                = "tapglue"
  password                = "demopasswd"
  # this should be true for production
  multi_az                = false
  publicly_accessible     = false
  vpc_security_group_ids  = [
    "${aws_security_group.rds_db.id}"]
  db_subnet_group_name    = "${aws_db_subnet_group.prod.id}"
  backup_retention_period = 7
  backup_window           = "04:00-04:30"
  maintenance_window      = "sat:05:00-sat:06:30"
}
/**/
# Database slaves
resource "aws_db_instance" "slave1" {
  identifier              = "slave1"
  storage_type            = "gp2"
  # change this to io1 if you want to use provisioned iops for production
  #iops = 3000 # this should give us a boost in performance for production
  allocated_storage       = "10"
  engine                  = "postgres"
  engine_version          = "9.4.4"
  instance_class          = "db.t2.micro"
  storage_encrypted       = false
  # if you want to change to true, see the list of instance types that support storage encryption: http://docs.aws.amazon.com/AmazonRDS/latest/UserGuide/Overview.Encryption.html#d0e10116
  name                    = "tapglue_prod"
  username                = "tapglue"
  password                = "demopasswd"
  multi_az                = false
  # this should be true for production
  publicly_accessible     = false
  replicate_source_db     = "${aws_db_instance.master.identifier}"
  vpc_security_group_ids  = [
    "${aws_security_group.rds_db.id}"]
  db_subnet_group_name    = "${aws_db_subnet_group.prod.id}"
  backup_retention_period = 0
  backup_window           = "04:00-04:30"
  maintenance_window      = "sat:05:00-sat:06:30"
}
/**/
