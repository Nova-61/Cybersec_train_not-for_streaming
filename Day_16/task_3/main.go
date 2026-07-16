package main

import (
	"fmt"
	"net/http"
)

func recoveryMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				fmt.Println("[PANIC] Тестовая паника!")
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
		}()
		next(w, r)
	}
}

func panicHandler(w http.ResponseWriter, r *http.Request) {
	panic("Тестовая паника!")
}

func main() {
	http.HandleFunc("/crash", recoveryMiddleware(panicHandler))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, World!")
	})

	fmt.Println("Server on http://localhost:8080")
	fmt.Println("  /panic - вызывает панику и перехватывает её")
	http.ListenAndServe(":8080", nil)
}