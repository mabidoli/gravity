# Gravity V2: AWS ECS Deployment

This directory contains all the necessary infrastructure-as-code and scripts to deploy the Gravity V2 application to AWS ECS Fargate.

## 1. Overview

The deployment is managed by:
- **Terraform**: For provisioning all AWS resources (VPC, ECS, RDS, Redis, ALB)
- **GitHub Actions**: For CI/CD, building Docker images, and deploying to ECS
- **Shell Scripts**: For local deployment and infrastructure setup

## 2. Prerequisites

Before you begin, you will need:
- An **AWS account** with administrative privileges
- **Terraform** (v1.5.0+) installed locally
- **AWS CLI** (v2+) installed and configured locally
- **Docker** installed locally
- A **GitHub repository** for your project

## 3. Setup Instructions

### Step 1: Configure AWS Credentials for GitHub Actions

1.  Create an IAM user in your AWS account with programmatic access.
2.  Attach the `AdministratorAccess` policy (for simplicity, or create a more restrictive policy later).
3.  In your GitHub repository, go to **Settings > Secrets and variables > Actions**.
4.  Create the following repository secrets:
    -   `AWS_ACCESS_KEY_ID`: The access key ID of the IAM user.
    -   `AWS_SECRET_ACCESS_KEY`: The secret access key of the IAM user.

### Step 2: Configure Terraform Backend (Optional but Recommended)

For team collaboration, it is highly recommended to use an S3 backend for Terraform state management.

1.  Create an S3 bucket in your AWS account (e.g., `gravity-terraform-state`).
2.  Create a DynamoDB table for state locking (e.g., `gravity-terraform-locks`) with a primary key named `LockID` (string).
3.  Uncomment and configure the `backend "s3"` block in `infrastructure/terraform/environments/prod/main.tf`.

### Step 3: Set Up Infrastructure

This step provisions all the necessary AWS resources using Terraform.

1.  Navigate to the environment directory:
    ```bash
    cd infrastructure/terraform/environments/prod
    ```
2.  Copy the example variables file:
    ```bash
    cp terraform.tfvars.example terraform.tfvars
    ```
3.  **Edit `terraform.tfvars`** and provide your specific configuration:
    -   `db_password`: A strong password for the RDS database.
    -   `frontend_image` / `backend_image`: Update with your AWS account ID.
    -   `certificate_arn` (optional): Add your ACM certificate ARN for HTTPS.

4.  Run the setup script from the project root:
    ```bash
    ./scripts/setup-infrastructure.sh prod
    ```
5.  Review the Terraform plan and type `yes` to apply the changes.

This will create the VPC, subnets, ALB, ECS cluster, RDS database, and Redis cache.

## 4. Deployment

There are two ways to deploy the application:

### Option A: Automated Deployment with GitHub Actions (Recommended)

-   Pushing to the `main` branch will automatically trigger the `Deploy to AWS ECS` workflow.
-   The workflow will build and push the Docker images to ECR and update the ECS services.
-   You can also trigger a manual deployment from the GitHub Actions tab.

### Option B: Manual Deployment from Local Machine

If you need to deploy manually, you can use the `deploy.sh` script.

1.  Ensure you are authenticated with AWS (`aws configure`).
2.  Run the script from the project root:
    ```bash
    ./scripts/deploy.sh prod
    ```

This script will:
1.  Log in to ECR.
2.  Build and push the frontend and backend Docker images.
3.  Update the ECS services to use the new images.
4.  Wait for the services to stabilize.

## 5. Accessing the Application

After a successful deployment, the Application Load Balancer DNS name will be printed in the Terraform output and the deployment script logs. You can access the application at `http://<alb_dns_name>`.

If you configured an ACM certificate and a custom domain, you can create a CNAME record in your DNS provider pointing your domain to the ALB DNS name.

## 6. Destroying the Infrastructure

To tear down all the AWS resources, run the following command from the environment directory:

```bash
cd infrastructure/terraform/environments/prod
terraform destroy
```

**Warning**: This will permanently delete all resources, including the database and its data. Use with caution.
