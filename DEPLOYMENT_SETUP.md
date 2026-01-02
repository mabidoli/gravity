# Gravity V2: AWS ECS Deployment Setup Complete ✅

All deployment infrastructure has been created and committed to the repository. Here's what you need to do to complete the setup:

## 📦 What Was Created

### Terraform Infrastructure (13 files)
- **VPC Module**: Complete network setup with public/private subnets, NAT gateways
- **ALB Module**: Application Load Balancer with routing rules for frontend/backend
- **ECS Module**: Fargate cluster with task definitions and services
- **RDS Module**: PostgreSQL 15 database with encryption and backups
- **Redis Module**: ElastiCache Redis for caching
- **Production Environment**: Complete Terraform configuration ready to deploy

### Deployment Scripts (2 files)
- `scripts/setup-infrastructure.sh`: Sets up AWS infrastructure with Terraform
- `scripts/deploy.sh`: Builds and deploys Docker images to ECS

### Docker Configuration
- `frontend/Dockerfile`: Multi-stage production build for Next.js
- `frontend/next.config.ts`: Updated with standalone output mode

### Documentation
- `infrastructure/README.md`: Complete deployment guide

---

## 🚀 Setup Instructions

### Step 1: Create AWS Account and Configure Credentials

1. **Create an AWS account** if you don't have one: https://aws.amazon.com/
2. **Install AWS CLI**: https://aws.amazon.com/cli/
3. **Configure AWS credentials**:
   ```bash
   aws configure
   ```
   Enter your:
   - AWS Access Key ID
   - AWS Secret Access Key
   - Default region (e.g., `us-east-1`)
   - Default output format (e.g., `json`)

### Step 2: Install Terraform

Download and install Terraform from: https://www.terraform.io/downloads

Verify installation:
```bash
terraform --version
```

### Step 3: Configure Terraform Variables

1. Navigate to the production environment:
   ```bash
   cd infrastructure/terraform/environments/prod
   ```

2. Copy the example variables file:
   ```bash
   cp terraform.tfvars.example terraform.tfvars
   ```

3. **Edit `terraform.tfvars`** with your configuration:
   ```hcl
   # AWS Configuration
   aws_region   = "us-east-1"  # Your preferred region
   project_name = "gravity"
   environment  = "prod"

   # Database Configuration
   # TODO: Generate a strong password
   db_password = "YOUR_STRONG_PASSWORD_HERE"

   # Docker Images
   # TODO: Update with your AWS account ID after running setup
   frontend_image = "123456789012.dkr.ecr.us-east-1.amazonaws.com/gravity-frontend:latest"
   backend_image  = "123456789012.dkr.ecr.us-east-1.amazonaws.com/gravity-backend:latest"

   # SSL Certificate (optional)
   # TODO: Create ACM certificate and add ARN here for HTTPS
   certificate_arn = ""
   ```

### Step 4: Set Up Infrastructure

From the project root, run:
```bash
./scripts/setup-infrastructure.sh prod
```

This script will:
1. Create ECR repositories for Docker images
2. Initialize Terraform
3. Plan infrastructure changes
4. Apply the infrastructure (after your confirmation)

**Expected time**: 10-15 minutes

**Cost estimate**: ~$50-100/month for production environment

### Step 5: Add GitHub Actions Workflow (Manual)

Due to GitHub App permissions, you need to manually add the CI/CD workflow:

1. Create the directory:
   ```bash
   mkdir -p .github/workflows
   ```

2. Copy the workflow file from the local repository:
   ```bash
   cp /path/to/gravity/.github/workflows/deploy-ecs.yml .github/workflows/
   ```

   Or create `.github/workflows/deploy-ecs.yml` with the content from the file at:
   `.github/workflows/deploy-ecs.yml` in your local repository.

3. **Add GitHub Secrets**:
   - Go to your GitHub repository
   - Navigate to **Settings > Secrets and variables > Actions**
   - Add these secrets:
     - `AWS_ACCESS_KEY_ID`: Your AWS access key
     - `AWS_SECRET_ACCESS_KEY`: Your AWS secret key

4. Commit and push:
   ```bash
   git add .github/workflows/deploy-ecs.yml
   git commit -m "ci: Add GitHub Actions deployment workflow"
   git push origin main
   ```

### Step 6: Deploy the Application

After infrastructure is set up, deploy the application:

```bash
./scripts/deploy.sh prod
```

This will:
1. Build Docker images for frontend and backend
2. Push images to ECR
3. Update ECS services
4. Wait for deployment to complete

---

## 🔑 Where to Add Your Keys

### AWS Credentials (Local Development)
```bash
aws configure
# Enter your keys when prompted
```

### GitHub Secrets (CI/CD)
Repository Settings > Secrets and variables > Actions:
- `AWS_ACCESS_KEY_ID`
- `AWS_SECRET_ACCESS_KEY`

### Terraform Variables
File: `infrastructure/terraform/environments/prod/terraform.tfvars`
- `db_password`: Database password
- `frontend_image`: Update with your AWS account ID
- `backend_image`: Update with your AWS account ID
- `certificate_arn`: (Optional) ACM certificate for HTTPS

---

## 📊 Infrastructure Overview

| Component | Service | Configuration |
|-----------|---------|---------------|
| **Compute** | ECS Fargate | 2 services (frontend + backend) |
| **Database** | RDS PostgreSQL 15 | db.t3.small, 50GB storage |
| **Cache** | ElastiCache Redis 7 | cache.t3.small |
| **Load Balancer** | Application LB | HTTP/HTTPS routing |
| **Network** | VPC | 2 AZs, public + private subnets |
| **Container Registry** | ECR | 2 repositories |

---

## 🎯 Next Steps

1. ✅ Create AWS account and configure credentials
2. ✅ Install Terraform
3. ⏳ Configure `terraform.tfvars` with your settings
4. ⏳ Run `./scripts/setup-infrastructure.sh prod`
5. ⏳ Add GitHub Actions workflow manually
6. ⏳ Run `./scripts/deploy.sh prod`
7. ⏳ Access your application at the ALB DNS name

---

## 📚 Documentation

- **Deployment Guide**: `infrastructure/README.md`
- **Backend Plans**: `backend/plans/`
- **Hosting Recommendation**: `backend/plans/06-hosting-recommendation.md`

---

## 💰 Cost Estimate

**Monthly costs for production environment**:
- ECS Fargate (2 services): ~$30-40
- RDS PostgreSQL (db.t3.small): ~$30-40
- ElastiCache Redis (cache.t3.small): ~$15-20
- Application Load Balancer: ~$20-25
- Data transfer & other: ~$10-15

**Total**: ~$105-140/month

**For development**: Use smaller instance types to reduce costs to ~$50-70/month

---

## 🆘 Support

If you encounter any issues:
1. Check the deployment logs
2. Review Terraform output for errors
3. Verify AWS credentials are configured correctly
4. Ensure all TODO items in configuration files are completed

Good luck with your deployment! 🚀
