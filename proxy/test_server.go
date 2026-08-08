package main

import (
    "encoding/json"
    "fmt"
    "net/http"
    "time"
)

type User struct {
    ID   int    `json:"id"`
    Name string `json:"name"`
}

var users = []User{
    {ID: 1, Name: "Ivan"},
    {ID: 2, Name: "Petr"},
    {ID: 3, Name: "Anna"},
}

func main() {
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        fmt.Printf("[SERVER] %s %s\n", r.Method, r.URL.Path)
        fmt.Fprintf(w, "Hello from target server!\n")
        fmt.Fprintf(w, "Path: %s\n", r.URL.Path)
        fmt.Fprintf(w, "Method: %s\n", r.Method)
    })

    http.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
        fmt.Printf("[SERVER] %s /users\n", r.Method)
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(users)
    })

    http.HandleFunc("/users/", func(w http.ResponseWriter, r *http.Request) {
        fmt.Printf("[SERVER] %s /users/{id}\n", r.Method)
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(users[0])
    })

    http.HandleFunc("/slow", func(w http.ResponseWriter, r *http.Request) {
        fmt.Printf("[SERVER] %s /slow\n", r.Method)
        time.Sleep(3 * time.Second)
        fmt.Fprintf(w, "Slow response after 3 seconds")
    })

    fmt.Println("Тестовый сервер запущен на :8080")
    fmt.Println("GET / - приветствие")
    fmt.Println("GET /users - список пользователей")
    fmt.Println("GET /users/1 - один пользователь")
    fmt.Println("GET /slow - медленный ответ (3 секунды)")
    http.ListenAndServe(":8080", nil)
}
