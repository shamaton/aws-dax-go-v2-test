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
  subnet_id              = aws_subnet.ec2.id
  vpc_security_group_ids = [aws_security_group.ec2.id]
  key_name               = aws_key_pair.this.key_name
  source_dest_check      = false
  iam_instance_profile   = aws_iam_instance_profile.ec2-profile.id

  tags = {
    Name = var.tag_name
  }
}

data "aws_iam_policy_document" "ec2-role" {

  statement {
    effect = "Allow"

    actions = ["sts:AssumeRole"]

    principals {
      type        = "Service"
      identifiers = ["ec2.amazonaws.com"]
    }
  }
}
# IAM Role for EC2 END

# IAM Role Policy for EC2
data "aws_iam_policy_document" "ec2-role_policy" {
  statement {
    effect = "Allow"
    actions = [
      "dynamodb:*",
      "dax:*",
    ]

    resources = [
      "*",
    ]
  }
}
# IAM Role Policy for EC2 END

# IAM Instance Profile
resource "aws_iam_instance_profile" "ec2-profile" {
  name = var.tag_name
  role = aws_iam_role.ec2-role.name
}
# IAM Instance Profile END

# Role
resource "aws_iam_role" "ec2-role" {
  name               = "${var.tag_name}-ec2-role"
  assume_role_policy = data.aws_iam_policy_document.ec2-role.json
}
# Role END

# Role-Policy
resource "aws_iam_role_policy" "ec2-role_policy" {
  name   = "${var.tag_name}-ec2-role-policy"
  role   = aws_iam_role.ec2-role.id
  policy = data.aws_iam_policy_document.ec2-role_policy.json
}


# SSH Command
output "command" {
  value = "./ssh.sh '${local_file.private_key_pem.filename}' '${aws_instance.this.public_ip}'"
}