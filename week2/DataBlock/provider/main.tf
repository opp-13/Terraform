provider "aws" {
  region = "ap-northeast-2"
}

provider "aws" {
  alias  = "us_east"
  region = "us-east-1"
}

data "aws_lambda_function" "lambda_ap" {
  provider      = aws
  function_name = "test_autoscale_up_lambda"
}

data "aws_lambda_function" "lambda_east" {
  provider      = aws.us_east
  function_name = "test_autoscale_up_lambda"
}

output "lambda_arns" {
  value = {
    ap_northeast_2 = data.aws_lambda_function.lambda_ap.arn
    us_east_1      = data.aws_lambda_function.lambda_east.arn
  }
}
