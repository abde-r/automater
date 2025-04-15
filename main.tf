provider "aws" {
    region = "us-east-1"
    access_key = var.access_key
    secret_key = var.secret_key
}

resource "aws_instance" "cloud-1" {
    ami = var.ami
    instance_type = "t2.micro"
    key_name = var.key_name
    security_groups = ["terraform_sg"]

    tags = {
        Name = "server-1"
    }
}

resource "aws_security_group" "terraform_sg" {
    name = "terraform_sg"
    description = "Allow SSH, HTTP, HTTPS, and all outbound traffic"

    ingress {
        description = "SSH"
        from_port   = 22
        to_port     = 22
        protocol    = "tcp"
        cidr_blocks = ["0.0.0.0/0"]
    }

    ingress {
        description = "TCP"
        from_port   = 8000
        to_port     = 9999
        protocol    = "tcp"
        cidr_blocks = ["0.0.0.0/0"]
    }

    egress {
        description = "Allow all outbound"
        from_port   = 0
        to_port     = 0
        protocol    = "-1"
        cidr_blocks = ["0.0.0.0/0"]
    }
}