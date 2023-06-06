
resource "aws_vpc" "vpc" {
  cidr_block           = var.vpc_cidr
  enable_dns_hostnames = true
  enable_dns_support   = true

  tags = {
    Name = var.tag_name
  }
}

resource "aws_subnet" "ec2" {
  vpc_id     = aws_vpc.vpc.id
  cidr_block = var.ec2_cidr
  availability_zone      = "us-east-1a"
  map_public_ip_on_launch = true

  tags = {
    Name = var.tag_name
  }
}

resource "aws_internet_gateway" "ec2" {
  vpc_id = aws_vpc.vpc.id

  tags = {
    Name = var.tag_name
  }
}

# Routing
resource "aws_route_table" "ec2" {
  vpc_id = aws_vpc.vpc.id

  tags = {
    Name = var.tag_name
  }
}

resource "aws_route" "ec2" {
  route_table_id         = aws_route_table.ec2.id
  gateway_id             = aws_internet_gateway.ec2.id
  destination_cidr_block = "0.0.0.0/0"
}

resource "aws_route_table_association" "ec2" {
  subnet_id      = aws_subnet.ec2.id
  route_table_id = aws_route_table.ec2.id
}

# Security Group
resource "aws_security_group" "ec2" {
  name        = "${var.tag_name}-ec2"
  vpc_id      = aws_vpc.vpc.id

  tags = {
    Name = var.tag_name
  }
}

resource "aws_security_group_rule" "ec2-ingress" {
  type              = "ingress"
  from_port         = 22
  to_port           = 22
  protocol          = "tcp"
  cidr_blocks       = var.ssh_cidr
  security_group_id = aws_security_group.ec2.id
}

resource "aws_security_group_rule" "ec2-egress" {
  type              = "egress"
  from_port         = 0
  to_port           = 0
  protocol          = "-1"
  cidr_blocks       = ["0.0.0.0/0"]
  security_group_id = aws_security_group.ec2.id
}


resource "aws_subnet" "dax" {
  vpc_id     = aws_vpc.vpc.id
  cidr_block = var.dax_cidr
  availability_zone      = "us-east-1a"

  tags = {
    Name = var.tag_name
  }
}

#resource "aws_security_group" "dax" {
#  name        = "${var.tag_name}-dax"
#  vpc_id      = aws_vpc.vpc.id
#
#  tags = {
#    Name = var.tag_name
#  }
#}

resource "aws_security_group_rule" "dax-ingress" {
  type              = "ingress"
  from_port         = 8111
  to_port           = 8111
  protocol          = "tcp"
  cidr_blocks       = [aws_vpc.vpc.cidr_block]
  security_group_id = aws_security_group.ec2.id
}