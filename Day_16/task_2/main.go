package main

import (
	"fmt"
	"net/http"
	"time"
)

func loggerMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		fmt.Printf("[LOG] Запрос к %s\n", r.URL.Path)
		next(w, r)
		duration := time.Since(start)
		fmt.Printf("[LOG] Время: %v\n", duration)
	}
}

func authMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		apiKey := r.Header.Get("X-API-Key")
		if apiKey != "123" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		next(w, r)
	}
}

func adminHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Admin panel")
}

func main() {
	http.HandleFunc("/admin", loggerMiddleware(authMiddleware(adminHandler)))

	fmt.Println("Server on http://localhost:8080")
	fmt.Println("  /admin - требует X-API-Key: 123")
	http.ListenAndServe(":8080", nil)
}