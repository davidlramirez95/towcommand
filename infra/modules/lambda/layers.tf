resource "aws_lambda_layer_version" "shared" {
  filename            = var.shared_layer_zip
  layer_name          = "towcommand-shared-${var.environment}"
  compatible_runtimes = ["nodejs.20.x"]
  compatible_architectures = ["arm64"]
  source_code_hash    = filebase64sha256(var.shared_layer_zip)

  tags = var.tags
}
