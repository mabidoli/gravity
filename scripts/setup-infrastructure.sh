#!/bin/bash

# Gravity V2 Infrastructure Setup Script
# This script sets up the initial AWS infrastructure using Terraform

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

print_info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

print_warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

print_step() {
    echo -e "${BLUE}[STEP]${NC} $1"
}

# Check if environment is provided
if [ -z "$1" ]; then
    print_error "Usage: ./setup-infrastructure.sh <environment>"
    print_info "Available environments: dev, staging, prod"
    exit 1
fi

ENVIRONMENT=$1
TERRAFORM_DIR="infrastructure/terraform/environments/$ENVIRONMENT"

print_info "Setting up infrastructure for environment: $ENVIRONMENT"

# Check if Terraform is installed
if ! command -v terraform &> /dev/null; then
    print_error "Terraform is not installed. Please install it first."
    print_info "Visit: https://www.terraform.io/downloads"
    exit 1
fi

# Check if AWS CLI is installed
if ! command -v aws &> /dev/null; then
    print_error "AWS CLI is not installed. Please install it first."
    exit 1
fi

# Check if user is logged in to AWS
if ! aws sts get-caller-identity &> /dev/null; then
    print_error "Not authenticated with AWS. Please run 'aws configure' first."
    exit 1
fi

# Get AWS account ID and region
AWS_ACCOUNT_ID=$(aws sts get-caller-identity --query Account --output text)
AWS_REGION=$(aws configure get region)
print_info "AWS Account ID: $AWS_ACCOUNT_ID"
print_info "AWS Region: $AWS_REGION"

# Check if terraform.tfvars exists
if [ ! -f "$TERRAFORM_DIR/terraform.tfvars" ]; then
    print_warn "terraform.tfvars not found. Creating from example..."
    cp "$TERRAFORM_DIR/terraform.tfvars.example" "$TERRAFORM_DIR/terraform.tfvars"
    print_error "Please edit $TERRAFORM_DIR/terraform.tfvars with your configuration before continuing."
    print_info "Required changes:"
    print_info "  1. Set a strong database password"
    print_info "  2. Update ECR repository URLs"
    print_info "  3. (Optional) Add ACM certificate ARN for HTTPS"
    exit 1
fi

# Create ECR repositories if they don't exist
print_step "Creating ECR repositories..."
aws ecr describe-repositories --repository-names gravity-frontend --region $AWS_REGION 2>/dev/null || \
    aws ecr create-repository --repository-name gravity-frontend --region $AWS_REGION > /dev/null
aws ecr describe-repositories --repository-names gravity-backend --region $AWS_REGION 2>/dev/null || \
    aws ecr create-repository --repository-name gravity-backend --region $AWS_REGION > /dev/null
print_info "✅ ECR repositories created"

# Navigate to Terraform directory
cd $TERRAFORM_DIR

# Initialize Terraform
print_step "Initializing Terraform..."
terraform init

# Validate configuration
print_step "Validating Terraform configuration..."
terraform validate

# Plan infrastructure
print_step "Planning infrastructure changes..."
terraform plan -out=tfplan

# Ask for confirmation
echo ""
print_warn "Review the plan above. Do you want to apply these changes?"
read -p "Type 'yes' to continue: " CONFIRM

if [ "$CONFIRM" != "yes" ]; then
    print_error "Deployment cancelled."
    exit 1
fi

# Apply infrastructure
print_step "Applying infrastructure changes..."
terraform apply tfplan

# Get outputs
print_info "Infrastructure deployment completed!"
echo ""
print_info "📋 Infrastructure Details:"
terraform output

# Run database migrations
print_step "Would you like to run database migrations now?"
read -p "Type 'yes' to run migrations: " RUN_MIGRATIONS

if [ "$RUN_MIGRATIONS" == "yes" ]; then
    print_info "Running database migrations..."
    # TODO: Add migration command here
    # For now, provide instructions
    print_warn "Please run migrations manually:"
    print_info "  1. Get the database endpoint from Terraform outputs"
    print_info "  2. Connect to the database"
    print_info "  3. Run: psql -h <endpoint> -U gravity_admin -d gravity_db -f backend/gravity-bff/migrations/0001_initial_schema.up.sql"
fi

print_info "✅ Infrastructure setup complete!"
print_info "Next steps:"
print_info "  1. Build and push Docker images"
print_info "  2. Run: ./scripts/deploy.sh $ENVIRONMENT"
