resource "aws_security_group" "ec-redis" {
  vpc_id      = "${aws_vpc.prod.id}"
  name        = "ec-redis"
  description = "Redis cache EC2 incoming traffic"

  ingress {
    from_port   = 6379
    to_port     = 6379
    protocol    = "tcp"
    cidr_blocks = [
      "${aws_subnet.frontend-a.cidr_block}"]
  }
  ingress {
    from_port   = 6379
    to_port     = 6379
    protocol    = "tcp"
    cidr_blocks = [
      "${aws_subnet.frontend-b.cidr_block}"]
  }

  ingress {
    from_port   = 6379
    to_port     = 6379
    protocol    = "tcp"
    cidr_blocks = [
      "${aws_subnet.backend-a.cidr_block}"]
  }
  ingress {
    from_port   = 6379
    to_port     = 6379
    protocol    = "tcp"
    cidr_blocks = [
      "${aws_subnet.backend-b.cidr_block}"]
  }

  egress {
    from_port   = 6379
    to_port     = 6379
    protocol    = "tcp"
    cidr_blocks = [
      "${aws_subnet.frontend-a.cidr_block}"]
  }
  egress {
    from_port   = 6379
    to_port     = 6379
    protocol    = "tcp"
    cidr_blocks = [
      "${aws_subnet.frontend-b.cidr_block}"]
  }

  egress {
    from_port   = 6379
    to_port     = 6379
    protocol    = "tcp"
    cidr_blocks = [
      "${aws_subnet.backend-a.cidr_block}"]
  }
  egress {
    from_port   = 6379
    to_port     = 6379
    protocol    = "tcp"
    cidr_blocks = [
      "${aws_subnet.backend-b.cidr_block}"]
  }

  tags {
    Name = "Redis cache incoming traffic"
  }
}

resource "aws_security_group" "ec-redis-ec2" {
  vpc_id      = "${aws_vpc.prod.id}"
  name        = "ec-redis-ec2"
  description = "Redis cache EC2 outgoing traffic"

  ingress {
    from_port       = 6379
    to_port         = 6379
    protocol        = "tcp"
    security_groups = [
      "${aws_security_group.ec-redis.id}"]
  }

  egress {
    from_port       = 6379
    to_port         = 6379
    protocol        = "tcp"
    security_groups = [
      "${aws_security_group.ec-redis.id}"]
  }

  tags {
    Name = "Redis cache incoming traffic"
  }
}

resource "aws_elasticache_subnet_group" "rate-limiter" {
  name        = "rate-limiter"
  description = "rate limiter cache group"
  subnet_ids  = [
    "${aws_subnet.frontend-a.id}",
    "${aws_subnet.frontend-b.id}"]
}

resource "aws_elasticache_cluster" "rate-limiter" {
  depends_on           = [
    "aws_elasticache_subnet_group.rate-limiter"]
  cluster_id           = "rate-limiter"
  engine               = "redis"
  engine_version       = "2.8.21"
  node_type            = "cache.t2.micro"
  port                 = 6379
  num_cache_nodes      = 1
  parameter_group_name = "default.redis2.8"
  maintenance_window   = "sun:05:00-sun:06:00"
  subnet_group_name    = "${aws_elasticache_subnet_group.rate-limiter.name}"
  security_group_ids   = [
    "${aws_security_group.ec-redis.id}"]
}
