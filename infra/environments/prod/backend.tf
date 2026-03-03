terraform {
  backend "s3" {
    bucket         = "towcommand-terraform-state-prod"
    key            = "infra/prod/terraform.tfstate"
    region         = "us-east-1"
    encrypt        = true
    dynamodb_table = "terraform-lock-prod"
  }
}
