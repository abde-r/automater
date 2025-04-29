package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

func main() {
	err := os.MkdirAll("terraform", 0755)
	if err != nil {
		fmt.Println("❌ Error creating terraform directory:", err)
		return
	}

	// Load environment variables from .env
	err = godotenv.Load(".env")
	if err != nil {
		fmt.Println("❌ Error loading .env file:", err)
		return
	}

	size := os.Getenv("DO_SIZE")
	token := os.Getenv("DO_TOKEN")
	image := os.Getenv("DO_IMAGE")
	region := os.Getenv("DO_REGION")
	sshKeyID := os.Getenv("DO_FINGER_PRINT")

	// Generate the Terraform file
	content := fmt.Sprintf(`
		terraform {
			required_providers {
				digitalocean = {
					source  = "digitalocean/digitalocean"
					version = "~> 2.0"
				}
			}
		}

		provider "digitalocean" {
			token = "%s"
		}

		resource "digitalocean_droplet" "cloud-1" {
			count = 2
			name   = "server"
			region = "%s"
			size   = "%s"
			image  = "%s"

			ssh_keys = ["%s"]

			tags = ["server"]
		}

		resource "digitalocean_firewall" "cloud_1_fw" {
			name        = "cloud-1-fw"
			droplet_ids = digitalocean_droplet.cloud-1[*].id

			inbound_rule {
				protocol         = "tcp"
				port_range       = "22"
				source_addresses = ["0.0.0.0/0", "::/0"]
			}

			inbound_rule {
				protocol         = "tcp"
				port_range       = "80"
				source_addresses = ["0.0.0.0/0", "::/0"]
			}

			inbound_rule {
				protocol         = "tcp"
				port_range       = "443"
				source_addresses = ["0.0.0.0/0", "::/0"]
			}

			inbound_rule {
				protocol         = "tcp"
				port_range       = "8180"
				source_addresses = ["0.0.0.0/0", "::/0"]
			}

			outbound_rule {
				protocol              = "tcp"
				port_range            = "all"
				destination_addresses = ["0.0.0.0/0", "::/0"]
			}

			outbound_rule {
				protocol              = "udp"
				port_range            = "53"
				destination_addresses = ["0.0.0.0/0", "::/0"]
			}

			outbound_rule {
				protocol              = "icmp"
				destination_addresses = ["0.0.0.0/0", "::/0"]
			}
		}
	`, token, region, size, image, sshKeyID)

	tfPath := filepath.Join("terraform", "main.tf")
	err = os.WriteFile(tfPath, []byte(content), 0644)
	if err != nil {
		fmt.Println("❌ Error writing main.tf:", err)
		return
	}
	fmt.Println("✅ Terraform file written: terraform/main.tf")
}