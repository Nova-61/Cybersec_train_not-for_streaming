package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Age   int    `json:"age"`
}

var users = []User{
	{ID: 1, Name: "Ivan", Email: "ivan@mail.com", Age: 30},
	{ID: 2, Name: "Petr", Email: "petr@mail.com", Age: 25},
}
var nextID = 3
var mu sync.Mutex

func getUsersHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

func getUserHandler(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	idStr := strings.TrimPrefix(path, "/users/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	for _, user := range users {
		if user.ID == id {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(user)
			return
		}
	}

	http.Error(w, "User not found", http.StatusNotFound)
}

func createUserHandler(w http.ResponseWriter, r *http.Request) {
	var newUser User
	err := json.NewDecoder(r.Body).Decode(&newUser)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if newUser.Name == "" || newUser.Email == "" {
		http.Error(w, "Name and Email are required", http.StatusBadRequest)
		return
	}

	if newUser.Age < 0 || newUser.Age > 150 {
		http.Error(w, "Age must be between 0 and 150", http.StatusBadRequest)
		return
	}

	mu.Lock()
	newUser.ID = nextID
	nextID++
	users = append(users, newUser)
	mu.Unlock()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newUser)
}

func deleteUserHandler(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	idStr := strings.TrimPrefix(path, "/users/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	mu.Lock()
	for i, user := range users {
		if user.ID == id {
			users = append(users[:i], users[i+1:]...)
			mu.Unlock()
			w.WriteHeader(http.StatusNoContent)
			return
		}
	}
	mu.Unlock()

	http.Error(w, "User not found", http.StatusNotFound)
}

func panicHandler(w http.ResponseWriter, r *http.Request) {
	panic("Тестовая паника!")
}

func recoveryMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				fmt.Printf("[PANIC] %v\n", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
		}()
		next(w, r)
	}
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "REST API Users\n")
		fmt.Fprintf(w, "GET /users - get all users\n")
		fmt.Fprintf(w, "GET /users/{id} - get user by ID\n")
		fmt.Fprintf(w, "POST /users - create user\n")
		fmt.Fprintf(w, "DELETE /users/{id} - delete user\n")
		fmt.Fprintf(w, "GET /panic - test panic recovery\n")
	})

	http.HandleFunc("/panic", recoveryMiddleware(panicHandler))
	http.HandleFunc("/users", recoveryMiddleware(getUsersHandler))

	http.HandleFunc("/users/", recoveryMiddleware(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			getUserHandler(w, r)
		case "POST":
			createUserHandler(w, r)
		case "DELETE":
			deleteUserHandler(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}))

	srv := &http.Server{
		Addr:    ":8080",
		Handler: nil,
	}

	go func() {
		fmt.Println("Server running on http://localhost:8080")
		fmt.Println("GET /users - get all users")
		fmt.Println("GET /users/1 - get user by ID")
		fmt.Println("POST /users - create user")
		fmt.Println("DELETE /users/1 - delete user")
		fmt.Println("GET /panic - test panic recovery")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("Server error: %v\n", err)
		}
	}()

	fmt.Println("\nPress Enter to shutdown...")
	fmt.Scanln()

	fmt.Println("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		fmt.Printf("Shutdown error: %v\n", err)
	}

	fmt.Println("Server gracefully stopped")
}
