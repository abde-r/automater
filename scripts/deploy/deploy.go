package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/joho/godotenv"
)

/*
	@writeInventoryFile
	write Inventory File
*/
func writeInventoryFile(ips []string, path string) error {
	var sb strings.Builder
	sb.WriteString("[droplets]\n")
	ansible_user := os.Getenv("ANSIBLE_USER")
	ansible_ssh_private_key_file := os.Getenv("ANSIBLE_SSH_PRIVATE_KEY_FILE")

	for _, ip := range ips {
		sb.WriteString(fmt.Sprintf("%s ansible_user=%s ansible_ssh_private_key_file=%s\n", ip, ansible_user, ansible_ssh_private_key_file))
	}

	return os.WriteFile(path, []byte(sb.String()), 0644)
}


/*
	@runAnsiblePlaybook
	Run Ansible Playbooks
*/
func runAnsiblePlaybook(playbookPath string, inventoryPath string) error {
	cmd := exec.Command("ansible-playbook", "-i", inventoryPath, playbookPath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	return cmd.Run()
}


/*
	@getInstanceIPsFromState
	Get Instance IPs From State
*/
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
		if !ok || resource["type"] != "digitalocean_droplet" {
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
			attr, ok := instance["attributes"].(map[string]interface{})
			if !ok {
				continue
			}

			ip, ok := attr["ipv4_address"].(string)
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


func main() {
	// Load environment variables from .env
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println("❌ Error loading .env file:", err)
		return
	}

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