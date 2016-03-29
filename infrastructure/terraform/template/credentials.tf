variable "cloudflare" {
  default = {
    email = "tools@tapglue.com"
    token = "8495c1d8eadae7413a79f74fa3bd3116c8c1b"
  }
  description = "CloudFlare credentials"
  type = "map"
}

variable "pg_db_name" {
  default = "tapglue"
}

variable "pg_username" {
  default = "tapglue"
}

variable "pg_password" {
  default = "gFthJy858iXOIA3IM0GARIuFYdIWkeCHJc0vto"
}
