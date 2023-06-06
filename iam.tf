# IAM for EC2
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

resource "aws_iam_instance_profile" "ec2-profile" {
  name = var.tag_name
  role = aws_iam_role.ec2-role.name
}

resource "aws_iam_role" "ec2-role" {
  name               = "${var.tag_name}-ec2-role"
  assume_role_policy = data.aws_iam_policy_document.ec2-role.json
}

resource "aws_iam_role_policy" "ec2-role_policy" {
  name   = "${var.tag_name}-ec2-role-policy"
  role   = aws_iam_role.ec2-role.id
  policy = data.aws_iam_policy_document.ec2-role_policy.json
}

# IAM for DAX / DynamoDB
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
