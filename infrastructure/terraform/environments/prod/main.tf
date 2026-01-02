# Production Environment Configuration for Gravity V2

terraform {
  required_version = ">= 1.5.0"

  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
  }

  # TODO: Configure S3 backend for state management
  # Uncomment and configure after creating S3 bucket
  # backend "s3" {
  #   bucket         = "gravity-terraform-state"
  #   key            = "prod/terraform.tfstate"
  #   region         = "us-east-1"
  #   encrypt        = true
  #   dynamodb_table = "gravity-terraform-locks"
  # }
}

provider "aws" {
  region = var.aws_region

  default_tags {
    tags = {
      Project     = "Gravity"
      Environment = "prod"
      ManagedBy   = "Terraform"
    }
  }
}

# Variables
variable "aws_region" {
  description = "AWS region"
  type        = string
  default     = "us-east-1"
}

variable "project_name" {
  description = "Project name"
  type        = string
  default     = "gravity"
}

variable "environment" {
  description = "Environment name"
  type        = string
  default     = "prod"
}

variable "db_password" {
  description = "Database master password"
  type        = string
  sensitive   = true
}

variable "frontend_image" {
  description = "Docker image for frontend"
  type        = string
  # TODO: Update with your ECR repository URL
  default     = "123456789012.dkr.ecr.us-east-1.amazonaws.com/gravity-frontend:latest"
}

variable "backend_image" {
  description = "Docker image for backend"
  type        = string
  # TODO: Update with your ECR repository URL
  default     = "123456789012.dkr.ecr.us-east-1.amazonaws.com/gravity-backend:latest"
}

variable "certificate_arn" {
  description = "ACM certificate ARN for HTTPS (optional)"
  type        = string
  default     = ""
}

# Data sources
data "aws_availability_zones" "available" {
  state = "available"
}

# VPC Module
module "vpc" {
  source = "../../modules/vpc"

  project_name       = var.project_name
  environment        = var.environment
  vpc_cidr           = "10.0.0.0/16"
  availability_zones = slice(data.aws_availability_zones.available.names, 0, 2)
}

# ALB Module
module "alb" {
  source = "../../modules/alb"

  project_name       = var.project_name
  environment        = var.environment
  vpc_id             = module.vpc.vpc_id
  public_subnet_ids  = module.vpc.public_subnet_ids
  certificate_arn    = var.certificate_arn
}

# RDS Module
module "rds" {
  source = "../../modules/rds"

  project_name           = var.project_name
  environment            = var.environment
  vpc_id                 = module.vpc.vpc_id
  private_subnet_ids     = module.vpc.private_subnet_ids
  ecs_security_group_id  = module.ecs.ecs_security_group_id
  db_password            = var.db_password
  db_instance_class      = "db.t3.small"
  allocated_storage      = 50
  max_allocated_storage  = 200
}

# Redis Module
module "redis" {
  source = "../../modules/redis"

  project_name           = var.project_name
  environment            = var.environment
  vpc_id                 = module.vpc.vpc_id
  private_subnet_ids     = module.vpc.private_subnet_ids
  ecs_security_group_id  = module.ecs.ecs_security_group_id
  node_type              = "cache.t3.small"
  num_cache_nodes        = 1
}

# ECS Module
module "ecs" {
  source = "../../modules/ecs"

  project_name               = var.project_name
  environment                = var.environment
  vpc_id                     = module.vpc.vpc_id
  private_subnet_ids         = module.vpc.private_subnet_ids
  alb_security_group_id      = module.alb.alb_security_group_id
  frontend_target_group_arn  = module.alb.frontend_target_group_arn
  backend_target_group_arn   = module.alb.backend_target_group_arn
  frontend_image             = var.frontend_image
  backend_image              = var.backend_image
  backend_db_host            = module.rds.db_address
  backend_db_name            = module.rds.db_name
  backend_db_user            = module.rds.db_username
  backend_db_password        = var.db_password
  backend_redis_host         = module.redis.redis_endpoint
}

# Outputs
output "alb_dns_name" {
  description = "DNS name of the load balancer"
  value       = module.alb.alb_dns_name
}

output "db_endpoint" {
  description = "RDS instance endpoint"
  value       = module.rds.db_endpoint
  sensitive   = true
}

output "redis_endpoint" {
  description = "Redis endpoint"
  value       = module.redis.redis_endpoint
  sensitive   = true
}

output "ecs_cluster_name" {
  description = "ECS Cluster name"
  value       = module.ecs.cluster_name
}
