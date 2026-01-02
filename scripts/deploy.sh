#!/bin/bash

# Gravity V2 Deployment Script
# This script helps deploy the application to AWS ECS

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Function to print colored output
print_info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

print_warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Check if environment is provided
if [ -z "$1" ]; then
    print_error "Usage: ./deploy.sh <environment>"
    print_info "Available environments: dev, staging, prod"
    exit 1
fi

ENVIRONMENT=$1
AWS_REGION=${AWS_REGION:-"us-east-1"}
PROJECT_NAME="gravity"

print_info "Deploying to environment: $ENVIRONMENT"
print_info "AWS Region: $AWS_REGION"

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

# Get AWS account ID
AWS_ACCOUNT_ID=$(aws sts get-caller-identity --query Account --output text)
print_info "AWS Account ID: $AWS_ACCOUNT_ID"

# ECR repository names
ECR_FRONTEND_REPO="$PROJECT_NAME-frontend"
ECR_BACKEND_REPO="$PROJECT_NAME-backend"

# ECR login
print_info "Logging in to Amazon ECR..."
aws ecr get-login-password --region $AWS_REGION | docker login --username AWS --password-stdin $AWS_ACCOUNT_ID.dkr.ecr.$AWS_REGION.amazonaws.com

# Build and push frontend
print_info "Building frontend Docker image..."
cd frontend
docker build -t $ECR_FRONTEND_REPO:latest .
docker tag $ECR_FRONTEND_REPO:latest $AWS_ACCOUNT_ID.dkr.ecr.$AWS_REGION.amazonaws.com/$ECR_FRONTEND_REPO:latest
print_info "Pushing frontend image to ECR..."
docker push $AWS_ACCOUNT_ID.dkr.ecr.$AWS_REGION.amazonaws.com/$ECR_FRONTEND_REPO:latest
cd ..

# Build and push backend
print_info "Building backend Docker image..."
cd backend/gravity-bff
docker build -t $ECR_BACKEND_REPO:latest .
docker tag $ECR_BACKEND_REPO:latest $AWS_ACCOUNT_ID.dkr.ecr.$AWS_REGION.amazonaws.com/$ECR_BACKEND_REPO:latest
print_info "Pushing backend image to ECR..."
docker push $AWS_ACCOUNT_ID.dkr.ecr.$AWS_REGION.amazonaws.com/$ECR_BACKEND_REPO:latest
cd ../..

# Update ECS services
CLUSTER_NAME="$PROJECT_NAME-$ENVIRONMENT-cluster"
FRONTEND_SERVICE="$PROJECT_NAME-$ENVIRONMENT-frontend-service"
BACKEND_SERVICE="$PROJECT_NAME-$ENVIRONMENT-backend-service"

print_info "Updating ECS services..."
aws ecs update-service --cluster $CLUSTER_NAME --service $FRONTEND_SERVICE --force-new-deployment --region $AWS_REGION > /dev/null
aws ecs update-service --cluster $CLUSTER_NAME --service $BACKEND_SERVICE --force-new-deployment --region $AWS_REGION > /dev/null

print_info "Waiting for services to stabilize..."
aws ecs wait services-stable --cluster $CLUSTER_NAME --services $FRONTEND_SERVICE --region $AWS_REGION
aws ecs wait services-stable --cluster $CLUSTER_NAME --services $BACKEND_SERVICE --region $AWS_REGION

print_info "✅ Deployment completed successfully!"
print_info "Environment: $ENVIRONMENT"
print_info "Cluster: $CLUSTER_NAME"

# Get ALB DNS name
ALB_DNS=$(aws elbv2 describe-load-balancers --region $AWS_REGION --query "LoadBalancers[?contains(LoadBalancerName, '$PROJECT_NAME-$ENVIRONMENT')].DNSName" --output text)
if [ -n "$ALB_DNS" ]; then
    print_info "Application URL: http://$ALB_DNS"
fi
