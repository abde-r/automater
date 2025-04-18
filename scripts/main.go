package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"github.com/joho/godotenv"
	"strings"
	"encoding/json"
)

/*
	Run Ansible Playbooks
*/
func runAnsiblePlaybook(playbookPath string, inventoryPath string) error {
	cmd := exec.Command("ansible-playbook", "-i", inventoryPath, playbookPath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	return cmd.Run()
}

func runCommand(args ...string) error {
	cmd := exec.Command("terraform", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	return cmd.Run()
}

func runTerraformApply() {
	err := os.Chdir("terraform")
	if err != nil {
		fmt.Println("❌ Failed to change directory:", err)
		return
	}

	fmt.Println("Initializing Terraform...")
	if err := runCommand("init"); err != nil {
		fmt.Println("❌ Error running terraform init:", err)
		return
	}

	fmt.Println("Planning Terraform...")
	if err := runCommand("plan"); err != nil {
		fmt.Println("❌ Error running terraform plan:", err)
		return
	}

	fmt.Println("Applying Terraform...")
	if err := runCommand("apply", "-auto-approve"); err != nil {
		fmt.Println("❌ Error running terraform apply:", err)
		return
	}

	fmt.Println("✅ Terraform apply completed successfully.")
}

func getInstanceIPsFromState(stateFilePath string) ([]string, error) {
	data, err := os.ReadFile(stateFilePath)
	if err != nil {
		return nil, fmt.Errorf("❌ failed to read state file: %w", err)
	}

	var tfState map[string]interface{}
	if err := json.Unmarshal(data, &tfState); err != nil {
		return nil, fmt.Errorf("❌ failed to parse JSON: %w", err)
	}

	resources, ok := tfState["resources"].([]interface{})
	if !ok {
		return nil, fmt.Errorf("❌ unexpected structure in state file: resources not found")
	}

	var ips []string
	for _, res := range resources {
		resource, ok := res.(map[string]interface{})
		if !ok || resource["type"] != "aws_instance" {
			continue
		}

		instances, ok := resource["instances"].([]interface{})
		if !ok {
			continue
		}

		for _, inst := range instances {
			instance, ok := inst.(map[string]interface{})
			if !ok {
				continue
			}
			attr := instance["attributes"].(map[string]interface{})
			ip, ok := attr["public_ip"].(string)
			if ok && ip != "" {
				ips = append(ips, ip)
			}
		}
	}

	if len(ips) == 0 {
		return nil, fmt.Errorf("no public IPs found")
	}

	return ips, nil
}


func writeInventoryFile(ips []string, path string) error {
	var sb strings.Builder
	sb.WriteString("[ec2]\n")
	ansible_ssh_private_key_file := os.Getenv("ANSIBLE_SSH_PRIVATE_KEY_FILE")
	ansible_user := os.Getenv("ANSIBLE_USER")

	for _, ip := range ips {
		sb.WriteString(fmt.Sprintf("%s ansible_user=%s ansible_ssh_private_key_file=%s\n", ip, ansible_user, ansible_ssh_private_key_file))
	}

	return os.WriteFile(path, []byte(sb.String()), 0644)
}

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

	accessKey := os.Getenv("AWS_ACCESS_KEY_ID")
	secretKey := os.Getenv("AWS_SECRET_ACCESS_KEY")
	ami := os.Getenv("AWS_AMI")
	keyName := os.Getenv("AWS_KEY_NAME")

	// Generate the Terraform file
	content := fmt.Sprintf(`
		provider "aws" {
			region     = "us-east-1"
			access_key = "%s"
			secret_key = "%s"
		}

		resource "aws_instance" "cloud-1" {
			ami           = "%s"
			instance_type = "t2.micro"
			key_name      = "%s"
			security_groups = ["terraform_sg"]
			count         = 2

			tags = {
				Name = "server"
			}
		}

		resource "aws_security_group" "terraform_sg" {
			name        = "terraform_sg"
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
	`, accessKey, secretKey, ami, keyName)

	tfPath := filepath.Join("terraform", "main.tf")
	err = os.WriteFile(tfPath, []byte(content), 0644)
	if err != nil {
		fmt.Println("❌ Error writing main.tf:", err)
		return
	}

	fmt.Println("✅ Terraform file written: terraform/main.tf")
	runTerraformApply()
	
	ips, err := getInstanceIPsFromState("terraform/terraform.tfstate")
	if err != nil {
		fmt.Println("❌ Error reading state:", err)
		return
	}

	if err := writeInventoryFile(ips, "ansible/inventory/inventory.ini"); err != nil {
		fmt.Println("❌ Failed to write inventory file:", err)
		return
	}

	fmt.Println("Running Ansible Setup Playbook...")
	if err := runAnsiblePlaybook("ansible/setup.yml", "ansible/inventory/inventory.ini"); err != nil {
		fmt.Println("❌ Error running Ansible Setup Playbook:", err)
		return
	}
	fmt.Println("Running Ansible Deploy Playbook...")
	if err := runAnsiblePlaybook("ansible/deploy.yml", "ansible/inventory/inventory.ini"); err != nil {
		fmt.Println("❌ Error running Ansible Deploy Playbook:", err)
		return
	}
	fmt.Println("✅ Ansible playbook completed successfully.")
}
