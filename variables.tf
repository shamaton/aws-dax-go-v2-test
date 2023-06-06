variable "tag_name" { default = "aws-dax-go-v2-test" }

variable "aws_profile" { default = "default" }

variable "ec2_instance_type" { default = "t3.micro" }

variable "vpc_cidr" {
  type = string
  default = "192.168.0.0/16"
}

variable "ec2_cidr" {
  type = string
  default = "192.168.1.0/24"
}

variable "dax_cidr" {
  type = string
  default = "192.168.2.0/24"
}

variable "ssh_cidr" {
  type = list(string)
  default = ["0.0.0.0/0"] # this setting is risky. recommend using tfvars.
}