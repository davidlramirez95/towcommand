#!/usr/bin/env bash
set -euo pipefail

STAGE=${1:-dev}
VALID_STAGES=("dev" "staging" "prod")

# Default to the requested AWS profile if none is set.
export AWS_PROFILE="${AWS_PROFILE:-iamadmin-general}"

if [[ ! " ${VALID_STAGES[*]} " =~ " ${STAGE} " ]]; then
  echo "❌ Invalid stage: $STAGE. Must be one of: ${VALID_STAGES[*]}"
  exit 1
fi

echo "🚀 Deploying TowCommand PH to '$STAGE'"
echo "========================================"

# Safety check for prod
if [ "$STAGE" = "prod" ]; then
  echo "⚠️  PRODUCTION deployment — are you sure? (y/N)"
  read -r confirmation
  if [ "$confirmation" != "y" ]; then
    echo "Deployment cancelled."
    exit 0
  fi
fi

# Run tests first
echo "🧪 Running tests..."
pnpm run test:unit

# Lint and typecheck
echo "🔍 Running lint and typecheck..."
pnpm run lint
pnpm run typecheck

# Build all packages
echo "🔨 Building packages..."
pnpm run build

# Prepare Terraform inputs (tfvars + Lambda ZIPs)
TF_ENV_DIR="infra/environments/${STAGE}"
TFVARS_PATH="${TF_ENV_DIR}/terraform.tfvars"
PKG_DIR="${TF_ENV_DIR}/packages"

mkdir -p "${PKG_DIR}"

if [[ ! -f "${TFVARS_PATH}" ]]; then
  echo "🧩 Creating ${TFVARS_PATH} from defaults (edit as needed)..."
  cat > "${TFVARS_PATH}" <<EOF
environment = "${STAGE}"
aws_region  = "us-east-1"

# VPC Configuration
vpc_cidr             = "10.0.0.0/16"
public_subnet_cidrs  = ["10.0.1.0/24", "10.0.2.0/24"]
private_subnet_cidrs = ["10.0.10.0/24", "10.0.11.0/24"]
availability_zones   = ["us-east-1a", "us-east-1b"]

# API Gateway
api_throttle_burst_limit        = 2000
api_throttle_rate_limit         = 1000
websocket_throttle_burst_limit  = 2000
websocket_throttle_rate_limit   = 1000

# Lambda
lambda_memory = 512

# Lambda Deployment Packages (paths to ZIP files)
booking_service_zip  = "./packages/booking-service.zip"
provider_service_zip = "./packages/provider-service.zip"
payment_service_zip  = "./packages/payment-service.zip"
sos_service_zip      = "./packages/sos-service.zip"
authorizer_zip       = "./packages/authorizer.zip"
shared_layer_zip     = "./packages/shared-layer.zip"

# Monitoring
alert_email        = "dev-alerts@example.com"
log_retention_days = 14

# Tags
tags = {
  CostCenter  = "Engineering"
  Owner       = "Platform Team"
  Application = "TowCommand"
}
EOF
fi

function make_node_lambda_zip() {
  local zip_path="$1"
  local label="$2"
  local tmp_dir
  tmp_dir="$(mktemp -d)"

  cat > "${tmp_dir}/index.js" <<'EOF'
export async function handler(event) {
  return {
    statusCode: 200,
    headers: { "content-type": "application/json" },
    body: JSON.stringify({
      ok: true,
      message: "TowCommand placeholder Lambda",
      receivedEventType: typeof event,
    }),
  };
}
EOF

  cat > "${tmp_dir}/package.json" <<'EOF'
{ "type": "module" }
EOF

  (cd "${tmp_dir}" && zip -qr "${OLDPWD}/${zip_path}" .)
  rm -rf "${tmp_dir}"
  echo "📦 Built ${label} -> ${zip_path}"
}

function make_node_layer_zip() {
  local zip_path="$1"
  local tmp_dir
  tmp_dir="$(mktemp -d)"
  mkdir -p "${tmp_dir}/nodejs"

  cat > "${tmp_dir}/nodejs/package.json" <<'EOF'
{ "name": "towcommand-shared-layer", "private": true, "type": "module" }
EOF

  (cd "${tmp_dir}" && zip -qr "${OLDPWD}/${zip_path}" .)
  rm -rf "${tmp_dir}"
  echo "📦 Built shared layer -> ${zip_path}"
}

make_node_lambda_zip "${PKG_DIR}/booking-service.zip" "booking-service"
make_node_lambda_zip "${PKG_DIR}/provider-service.zip" "provider-service"
make_node_lambda_zip "${PKG_DIR}/payment-service.zip" "payment-service"
make_node_lambda_zip "${PKG_DIR}/sos-service.zip" "sos-service"
make_node_lambda_zip "${PKG_DIR}/authorizer.zip" "authorizer"
make_node_layer_zip "${PKG_DIR}/shared-layer.zip"

# Deploy infrastructure
echo "🏗️ Deploying infrastructure..."
cd infra/environments/$STAGE
terraform init -upgrade
terraform plan -out=tfplan
terraform apply tfplan
cd ../../..

echo ""
echo "✅ Deployment to '$STAGE' complete!"
