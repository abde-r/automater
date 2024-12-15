package main

import (
	"fmt"
	"net/http"
	"os"
)

func homeHandler(w http.ResponseWriter, r *http.Request) {
	hostname, _ := os.Hostname()
	pageContent := `
	<!DOCTYPE html>
	<html lang="en">
	<head>
		<meta charset="UTF-8">
		<meta name="viewport" content="width=device-width, initial-scale=1.0">
		<title>Greeting App</title>
		<style>
			body {
				font-family: Arial, sans-serif;
				text-align: center;
				margin: 50px;
				background-color: #f4f4f9;
				color: #333;
			}
			.container {
				background-color: #fff;
				padding: 20px;
				border-radius: 10px;
				box-shadow: 0 4px 8px rgba(0, 0, 0, 0.2);
			}
			h1 {
				color: #007BFF;
			}
			p {
				font-size: 18px;
			}
			footer {
				margin-top: 20px;
				font-size: 14px;
				color: #777;
			}
		</style>
	</head>
	<body>
		<div class="container">
			<h1>Greeting App 1</h1>
			<p>Welcome to the Greeting Application deployed by Ansible!</p>
			<p>This application is running in a Kubernetes environment.</p>
			<p><strong>Pod Hostname:</strong> %s</p>
		</div>
		<footer>
			<p>&copy; 2024 Greeting App | Powered by Go and Kubernetes</p>
		</footer>
	</body>
	</html>
	`

	fmt.Fprintf(w, pageContent, hostname)
}

func main() {
	http.HandleFunc("/", homeHandler)

	port := "8080"
	fmt.Printf("Starting server on :%s\n", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		fmt.Println("Error starting server:", err)
	}
}