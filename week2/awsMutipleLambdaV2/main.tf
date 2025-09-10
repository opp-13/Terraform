terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 4.0"
    }
  }
}

provider "aws" {
  region = "ap-northeast-2"
}

# IAM Role for both Lambda functions
resource "aws_iam_role" "lambda_exec_role" {
  name = "lambda_exec_role"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [{
      Effect    = "Allow"
      Principal = { Service = "lambda.amazonaws.com" }
      Action    = "sts:AssumeRole"
    }]
  })
}

# Load the basic policy JSON for main Lambda
data "local_file" "basic_policy" {
  filename = "${path.module}/lambda_basic_policy.json"
}

# Load the full policy JSON for async Lambda
data "local_file" "full_policy" {
  filename = "${path.module}/lambda_full_policy.json"
}

# Create IAM Policy for main Lambda
resource "aws_iam_policy" "main_lambda_policy" {
  name   = "main_lambda_policy"
  policy = data.local_file.basic_policy.content
}

# Create IAM Policy for async Lambda
resource "aws_iam_policy" "async_lambda_policy" {
  name   = "async_lambda_policy"
  policy = data.local_file.full_policy.content
}

# Attach basic policy to IAM role (for main Lambda)
resource "aws_iam_role_policy_attachment" "main_lambda_attach" {
  role       = aws_iam_role.lambda_exec_role.name
  policy_arn = aws_iam_policy.main_lambda_policy.arn
}

# Attach full policy to IAM role (for async Lambda)
resource "aws_iam_role_policy_attachment" "async_lambda_attach" {
  role       = aws_iam_role.lambda_exec_role.name
  policy_arn = aws_iam_policy.async_lambda_policy.arn
}

# Package main Lambda code
data "archive_file" "main_lambda_zip" {
  type        = "zip"
  source_dir  = "${path.module}/main_lambda"
  output_path = "${path.module}/main_lambda.zip"
}

# Package async Lambda code
data "archive_file" "async_lambda_zip" {
  type        = "zip"
  source_dir  = "${path.module}/async_lambda"
  output_path = "${path.module}/async_lambda.zip"
}

# Deploy main Lambda function
resource "aws_lambda_function" "main_lambda" {
  function_name = "SlackCommandMainLambda"
  handler       = "lambda_function.lambda_handler"
  runtime       = "python3.10"
  role          = aws_iam_role.lambda_exec_role.arn
  filename      = data.archive_file.main_lambda_zip.output_path

  environment {
    variables = {
      ASYNC_LAMBDA_NAME = aws_lambda_function.async_lambda.function_name
    }
  }
}

# Create public Lambda Function URL for main handler
resource "aws_lambda_function_url" "main_lambda_url" {
  function_name      = aws_lambda_function.main_lambda.function_name
  authorization_type = "NONE"
  invoke_mode        = "BUFFERED"
}

# Deploy async Lambda function
resource "aws_lambda_function" "async_lambda" {
  function_name      = "SlackCommandAsyncLambda"
  handler            = "lambda_function.lambda_handler"
  runtime            = "python3.10"
  role               = aws_iam_role.lambda_exec_role.arn
  filename           = data.archive_file.async_lambda_zip.output_path
}

# Output main Lambda URL (Slack Slash Command 요청 URL)
output "main_lambda_url" {
  description = "Public URL of main Lambda function (Slack Slash Command endpoint)"
  value       = aws_lambda_function_url.main_lambda_url.function_url
}
