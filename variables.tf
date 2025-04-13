variable "access_key" {
  description = "aws user access key"
  type        = string
  sensitive   = true
}

variable "secret_key" {
  description = "aws user secret key"
  type        = string
  sensitive   = true
}

variable "ami" {
  description = "instance ami"
  type        = string
  sensitive   = true
}

variable "key_name" {
  description = "instance key name"
  type        = string
  sensitive   = true
}
