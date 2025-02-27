# EC2 Instance
resource "aws_instance" "web" {
  ami                    = var.ec2_ami
  instance_type          = var.ec2_instance_type
  subnet_id              = aws_subnet.public.id
  vpc_security_group_ids = [aws_security_group.ec2_sg.id]

  user_data = <<-EOF
              #!/bin/bash
              # Update and install packages
              apt-get update
              apt-get install -y apt-transport-https ca-certificates curl software-properties-common

              # Install Docker
              curl -fsSL https://download.docker.com/linux/ubuntu/gpg | apt-key add -
              add-apt-repository "deb [arch=amd64] https://download.docker.com/linux/ubuntu $(lsb_release -cs) stable"
              apt-get update
              apt-get install -y docker-ce docker-ce-cli containerd.io
              systemctl enable docker
              systemctl start docker
              usermod -aG docker ubuntu

              # Install Docker Compose
              DOCKER_COMPOSE_VERSION=$(curl -s https://api.github.com/repos/docker/compose/releases/latest | grep 'tag_name' | cut -d\" -f4)
              mkdir -p /usr/local/lib/docker/cli-plugins
              curl -L "https://github.com/docker/compose/releases/download/$${DOCKER_COMPOSE_VERSION}/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/lib/docker/cli-plugins/docker-compose
              chmod +x /usr/local/lib/docker/cli-plugins/docker-compose
              ln -s /usr/local/lib/docker/cli-plugins/docker-compose /usr/local/bin/docker-compose


              # Enable password authentication for SSH
              sed -i 's/PasswordAuthentication no/PasswordAuthentication yes/' /etc/ssh/sshd_config
              systemctl restart sshd

              # Set password for ubuntu user
              echo '${var.ec2_username}:${var.ec2_password}' | chpasswd
              EOF

  tags = {
    Name = "web-server"
    Project   = "test-app"
  }
}