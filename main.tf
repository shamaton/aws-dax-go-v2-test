terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 3.0"
    }
  }
}

# Configure the AWS Provider
provider "aws" {
  region  = "us-east-1"
#  access_key = var.aws_access_key
#  secret_key = var.aws_secret_key
  profile = var.aws_profile
}

# Get latest Amazon Linux 2 AMI
data "aws_ami" "amazon-linux-2" {
  most_recent = true
  owners      = ["amazon"]
  filter {
    name   = "name"
    values = ["amzn2-ami-hvm*"]
  }
}

# Key Pair for SSH to EC2 Instance
resource "aws_key_pair" "this" {
  key_name   = var.key_name
  public_key = tls_private_key.keygen.public_key_openssh
}

# EC2 Instance
resource "aws_instance" "this" {
  ami                    = data.aws_ami.amazon-linux-2.id
  instance_type          = var.ec2_instance_type
  subnet_id              = aws_subnet.this.id
  vpc_security_group_ids = [aws_security_group.this.id]
  key_name               = aws_key_pair.this.key_name
  source_dest_check      = false

  tags = {
    Name = var.tag_name
  }
}

# SSH Command
output "command" {
  value = "ssh -i ${local_file.private_key_pem.filename} ec2-user@${aws_instance.this.public_ip}"
}