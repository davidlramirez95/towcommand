terraform {
  backend "s3" {
    bucket         = "towcommand-terraform-state-staging"
    key            = "infra/staging/terraform.tfstate"
    region         = "us-east-1"
    encrypt        = true
    dynamodb_table = "terraform-lock-staging"
  }
}
