package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/joho/godotenv"
)


func runCommand(args ...string) error {
	cmd := exec.Command("terraform", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	return cmd.Run()
}


/*
	@runTerraformApply
	run Terraform Apply
*/
func runTerraformApply() {
    // 1. Save original working directory
    oldWd, err := os.Getwd()
    if err != nil {
        log.Fatalf("❌ unable to get current dir: %v", err)
    }
    // 2. Change into the terraform directory
    if err := os.Chdir("terraform"); err != nil {
        log.Fatalf("❌ failed to change directory: %v", err)
    }
    // 3. Ensure we restore the original directory
    defer func() {
        if err := os.Chdir(oldWd); err != nil {
            log.Fatalf("❌ failed to restore directory: %v", err)
        }
    }()

    fmt.Println("Initializing Terraform...")
    if err := runCommand("init"); err != nil {
        log.Fatalf("❌ terraform init failed: %v", err)
    }

    fmt.Println("Planning Terraform...")
    if err := runCommand("plan"); err != nil {
        log.Fatalf("❌ terraform plan failed: %v", err)
    }

    fmt.Println("Applying Terraform...")
    if err := runCommand("apply", "-auto-approve"); err != nil {
        log.Fatalf("❌ terraform apply failed: %v", err)
    }

    fmt.Println("✅ Terraform apply completed successfully.")
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

	runTerraformApply()
	fmt.Println("✅ All instances are ready.")
}