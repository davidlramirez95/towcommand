terraform {
  backend "s3" {
    bucket         = "towcommand-terraform-state-dev"
    key            = "infra/dev/terraform.tfstate"
    region         = "us-east-1"
    encrypt        = true
    dynamodb_table = "terraform-lock-dev"
  }
}
