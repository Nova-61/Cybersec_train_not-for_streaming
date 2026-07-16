package main

import (
	"fmt"
	"net/http"
	"sync"
)

type UserHandler struct {
	Username string
}

func (h UserHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, %s!, ты ЧМО", h.Username)
}

type CountHandler struct {
	Count int
	mu    sync.Mutex
}

func (h *CountHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.mu.Lock()
	h.Count++
	h.mu.Unlock()

	fmt.Fprintf(w, "Колличество запросов: %d", h.Count)
}

func main() {
	userHandler := &UserHandler{Username: "Тима"}
	countHandler := &CountHandler{}

	http.Handle("/user", userHandler)
	http.Handle("/counter", countHandler)

	fmt.Println("Server running on http://localhost:8080")
	fmt.Println("  /user    - Hello, John!")
	fmt.Println("  /counter - Количество запросов")

	http.ListenAndServe(":8080", nil)
}
