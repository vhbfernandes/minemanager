provider "aws" {
    region = "us-east-1"
}

#preencher
locals {
  user_data = <<EOF
#!/bin/bash
wget -O factorio_headless.tar.gz https://www.factorio.com/get-download/1.1.53/headless/linux64
tar -xf ~/factorio_headless.tar.gz
mkdir factorio/saves
/root/factorio/bin/x64/factorio --start-server /root/factorio/saves/thegame.zip
EOF
}

#preencher
resource "aws_key_pair" "accesskey" {
  key_name   = "access-key"
  public_key = "" #your pubkey
}

resource "aws_instance" "aws_instance_suffix" {
  ami = data.aws_ami.ubuntu.id
  instance_type = "t2.small"
  key_name = aws_key_pair.accesskey.key_name
  user_data_base64 = base64encode(local.user_data)
  vpc_security_group_ids = [aws_security_group.allow_ssh.id]
  tags = {
    "instance" = "factorio"
  }
}

#raramente vai precisar mexer

data "aws_ami" "ubuntu" {
  most_recent = true

  filter {
    name   = "name"
    values = ["ubuntu/images/hvm-ssd/ubuntu-focal-20.04-amd64-server-*"]
  }

  filter {
    name   = "virtualization-type"
    values = ["hvm"]
  }

  owners = ["099720109477"] # Canonical
}

resource "aws_security_group" "allow_ssh" {
  name        = "allow_ssh"
  description = "Allow ssh traffic"
  vpc_id      = data.aws_vpc.default.id

  ingress {
    description      = "SSH"
    from_port        = 22
    to_port          = 22
    protocol         = "tcp"
    cidr_blocks      = ["0.0.0.0/0"]
    ipv6_cidr_blocks = ["::/0"]
  }
  ingress {
    description = "FACTORIO"
    from_port = 34197
    to_port = 34197
    protocol = "udp"
    cidr_blocks      = ["0.0.0.0/0"]
    ipv6_cidr_blocks = ["::/0"]
  }
  egress {
    from_port        = 0
    to_port          = 0
    protocol         = "-1"
    cidr_blocks      = ["0.0.0.0/0"]
    ipv6_cidr_blocks = ["::/0"]
  }
}


data "aws_vpc" "default" {
  default = true
}