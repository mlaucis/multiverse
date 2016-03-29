variable "ami" {
  default     = {}
  description = "AMIs used for components"
  type        = "map"
}
variable "env" {
  default     = ""
  description = "Environment name used for isolation"
  type        = "string"
}

variable "region" {
  default     = ""
  description = "Region to deploy to"
  type        = "string"
}

variable "role" {
  default     = {
    "rds-monitoring-role" = "arn:aws:iam::775034650473:role/rds-monitoring-role"
  }
  description = "Roles shared between envs"
  type        = "map"
}

variable "version" {
  default     = {}
  description = "Versions used for deployed services"
  type        = "map"
}
