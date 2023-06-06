resource "aws_dynamodb_table" "this" {
  name           = "GameScores"
  billing_mode   = "PROVISIONED"
  read_capacity  = 20
  write_capacity = 20
  hash_key       = "UserId"
  range_key      = "GameTitle"

  attribute {
    name = "UserId"
    type = "S"
  }

  attribute {
    name = "GameTitle"
    type = "S"
  }

  attribute {
    name = "TopScore"
    type = "N"
  }

  ttl {
    attribute_name = "TimeToExist"
    enabled        = false
  }

  global_secondary_index {
    name               = "GameTitleIndex"
    hash_key           = "GameTitle"
    range_key          = "TopScore"
    write_capacity     = 10
    read_capacity      = 10
    projection_type    = "INCLUDE"
    non_key_attributes = ["UserId"]
  }

  tags = {
    Name        = var.tag_name
    Environment = "test"
  }
}

resource "aws_dax_cluster" "this" {
  cluster_name         = var.tag_name
  iam_role_arn         = aws_iam_role.dax.arn
  node_type            = "dax.t3.small"
  replication_factor   = 1
  security_group_ids   = [aws_security_group.ec2.id]
  subnet_group_name    = aws_dax_subnet_group.this.id
  parameter_group_name = aws_dax_parameter_group.this.id
}

resource "aws_dax_subnet_group" "this" {
  name       = var.tag_name
  subnet_ids = [aws_subnet.dax.id]
}

resource "aws_dax_parameter_group" "this" {
  name = var.tag_name

  parameters {
    name  = "query-ttl-millis"
    value = "100000"
  }

  parameters {
    name  = "record-ttl-millis"
    value = "100000"
  }
}

resource "aws_iam_role" "dax" {
  name               = var.tag_name
  assume_role_policy = data.aws_iam_policy_document.dax_assume_role.json
}

resource "aws_iam_role_policy" "dax" {
  role   = aws_iam_role.dax.id
  policy = data.aws_iam_policy_document.dax.json
}

data "aws_iam_policy_document" "dax_assume_role" {
  statement {
    actions = ["sts:AssumeRole"]
    principals {
      type        = "Service"
      identifiers = ["dax.amazonaws.com"]
    }
  }
}

data "aws_iam_policy_document" "dax" {
  statement {
    sid = "DynamoDBAccess"
    actions = [
      "dynamodb:DescribeTable",
      "dynamodb:PutItem",
      "dynamodb:GetItem",
      "dynamodb:UpdateItem",
      "dynamodb:DeleteItem",
      "dynamodb:Query",
      "dynamodb:Scan",
      "dynamodb:BatchGetItem",
      "dynamodb:BatchWriteItem",
      "dynamodb:ConditionCheckItem",
    ]
    resources = [
      aws_dynamodb_table.this.arn,
      "${aws_dynamodb_table.this.arn}/index/*"
    ]
  }
}

## Security Group
#resource "aws_security_group" "dax" {
#  name        = var.tag_name
#  vpc_id      = aws_vpc.vpc.id
#
#  tags = {
#    Name = var.tag_name
#  }
#}
#
#resource "aws_security_group_rule" "dax-ingress" {
#  type              = "ingress"
#  from_port         = 8111
#  to_port           = 8111
#  protocol          = "tcp"
#  cidr_blocks       = ["0.0.0.0/0"]
#  security_group_id = aws_security_group.dax.id
#}
#
#resource "aws_security_group_rule" "dax-egress" {
#  type              = "egress"
#  from_port         = 0
#  to_port           = 0
#  protocol          = "-1"
#  cidr_blocks       = ["0.0.0.0/0"]
#  security_group_id = aws_security_group.dax.id
#}