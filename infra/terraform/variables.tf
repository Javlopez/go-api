variable "aws_region" {
  description = "AWS region"
  type        = string
  default     = "us-east-1"
}

variable "aws_shared_config_files" {
  description = "List of aws shared config files"
  type        = list(string)
}

variable "aws_shared_credentials_files" {
  description = "List of aws shared credentials files"
  type        = list(string)
}

variable "aws_profile" {
  description = "Name of the profile, empty if it is the unique"
  type = string
}


variable "vpc_cidr" {
  description = "CIDR block for the VPC"
  type        = string
  default     = "10.0.0.0/16"
}

variable "public_subnet_cidr" {
  description = "CIDR block for the public subnet"
  type        = string
  default     = "10.0.1.0/24"
}

variable "availability_zone_public" {
  description = "Availability zone for the public subnet"
  type        = string
  default     = "us-east-1a"
}

variable "ec2_instance_type" {
  description = "EC2 instance type"
  type        = string
  default     = "t2.small"
}

variable "ec2_ami" {
  description = "AMI ID for EC2 instance"
  type        = string
  default     = "ami-0261755bbcb8c4a84" # Ubuntu 20.04 LTS in us-east-1
}

variable "ec2_username" {
  description = "Username for EC2 instance SSH access"
  type        = string
  default     = "ubuntu"
}

variable "ec2_password" {
  description = "Password for EC2 instance SSH access"
  type        = string
  default     = "StrongPassword123!" # Change this to a secure password
  sensitive   = true
}