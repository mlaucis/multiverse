resource "aws_elasticache_subnet_group" "rate-limiter" {
  name        = "rate-limiter"
  description = "rate limiter cache group"
  subnet_ids  = [
    "${aws_subnet.frontend-a.id}",
    "${aws_subnet.frontend-b.id}"]
}

resource "aws_elasticache_cluster" "rate-limiter" {
  depends_on           = [
    "aws_elasticache_subnet_group.rate-limiter",
  ]
  cluster_id           = "rate-limiter"
  engine               = "redis"
  engine_version       = "2.8.21"
  node_type            = "cache.r3.large"
  port                 = 6379
  num_cache_nodes      = 1
  parameter_group_name = "default.redis2.8"
  maintenance_window   = "sun:05:00-sun:06:00"
  subnet_group_name    = "${aws_elasticache_subnet_group.rate-limiter.name}"
  security_group_ids   = [
    "${aws_security_group.platform.id}",
  ]
}
