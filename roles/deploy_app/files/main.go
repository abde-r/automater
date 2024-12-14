package main

import (
    "fmt"
    "net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Welcome to the Go Server deployed by Ansible!")
}

func main() {
    http.HandleFunc("/", handler)
    fmt.Println("Starting Go server on :8080")
    http.ListenAndServe(":8080", nil)
}
