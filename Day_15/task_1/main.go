package main

import (
	"fmt"
	"net/http"
)

func helloHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, World!")
}

func pingHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "pong")
}

func infoHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Server: Go, Version: 1.22")
}

func main() {
	http.HandleFunc("/hello", helloHandler)
	http.HandleFunc("/ping", pingHandler)
	http.HandleFunc("/info", infoHandler)

	fmt.Println("Server running on http://localhost:8080")
	fmt.Println("Available endpoints:")
	fmt.Println("  /hello - Hello, World!")
	fmt.Println("  /ping  - pong")
	fmt.Println("  /info  - Server info")

	http.ListenAndServe(":8080", nil)
}
