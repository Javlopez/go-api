provider "aws" {
  region  = var.aws_region
  shared_config_files      = var.aws_shared_config_files
  shared_credentials_files = var.aws_shared_credentials_files
  profile                  = var.aws_profile
}

# All resources are defined in their respective module files:
# - network.tf: VPC, subnets, route tables, security groups
# - ec2.tf: EC2 instance with Docker and Docker Compose
# - outputs.tf: Output values