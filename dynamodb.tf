resource "aws_dynamodb_table" "this" {
  name           = "GameScores"
  billing_mode   = "PROVISIONED"
  read_capacity  = 20
  write_capacity = 20
  hash_key       = "user_id"
  range_key      = "game_title"

  attribute {
    name = "user_id"
    type = "S"
  }

  attribute {
    name = "game_title"
    type = "S"
  }

  attribute {
    name = "score"
    type = "N"
  }

  ttl {
    attribute_name = "TimeToExist"
    enabled        = false
  }

  global_secondary_index {
    name               = "GameTitleIndex"
    hash_key           = "game_title"
    range_key          = "score"
    write_capacity     = 10
    read_capacity      = 10
    projection_type    = "INCLUDE"
    non_key_attributes = ["user_id"]
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