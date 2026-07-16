package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"sync"
)

// Структура пользователя
type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Age   int    `json:"age"`
}

// Хранилище пользователей (в памяти)
var users = []User{
	{ID: 1, Name: "Ivan", Email: "ivan@mail.com", Age: 30},
	{ID: 2, Name: "Petr", Email: "petr@mail.com", Age: 25},
}
var nextID = 3    // Следующий ID для нового пользователя
var mu sync.Mutex // Для защиты от гонки данных

// --- ОБРАБОТЧИКИ ---

// 1. GET /users — получить всех пользователей
func getUsersHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

// 2. GET /users/{id} — получить одного пользователя
func getUserHandler(w http.ResponseWriter, r *http.Request) {
	// Получаем ID из пути: /users/1 → id = 1
	path := r.URL.Path
	idStr := strings.TrimPrefix(path, "/users/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	// Ищем пользователя
	for _, user := range users {
		if user.ID == id {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(user)
			return
		}
	}

	http.Error(w, "User not found", http.StatusNotFound)
}

// 3. POST /users — создать пользователя
func createUserHandler(w http.ResponseWriter, r *http.Request) {
	// Читаем JSON из тела запроса
	var newUser User
	err := json.NewDecoder(r.Body).Decode(&newUser)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Проверяем обязательные поля
	if newUser.Name == "" || newUser.Email == "" {
		http.Error(w, "Name and Email are required", http.StatusBadRequest)
		return
	}

	// Добавляем пользователя с защитой от гонки
	mu.Lock()
	newUser.ID = nextID
	nextID++
	users = append(users, newUser)
	mu.Unlock()

	// Возвращаем созданного пользователя
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newUser)
}

// 4. DELETE /users/{id} — удалить пользователя
func deleteUserHandler(w http.ResponseWriter, r *http.Request) {
	// Получаем ID из пути
	path := r.URL.Path
	idStr := strings.TrimPrefix(path, "/users/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	// Ищем и удаляем пользователя
	mu.Lock()
	for i, user := range users {
		if user.ID == id {
			// Удаляем элемент из слайса
			users = append(users[:i], users[i+1:]...)
			mu.Unlock()
			w.WriteHeader(http.StatusNoContent)
			return
		}
	}
	mu.Unlock()

	http.Error(w, "User not found", http.StatusNotFound)
}

// --- ГЛАВНАЯ ФУНКЦИЯ ---

func main() {
	// Главная страница (информация)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "REST API Users\n")
		fmt.Fprintf(w, "GET /users - get all users\n")
		fmt.Fprintf(w, "GET /users/{id} - get user by ID\n")
		fmt.Fprintf(w, "POST /users - create user\n")
		fmt.Fprintf(w, "DELETE /users/{id} - delete user\n")
	})

	// GET /users — все пользователи
	http.HandleFunc("/users", getUsersHandler)

	// GET /users/{id} — один пользователь
	// POST /users — создать
	// DELETE /users/{id} — удалить
	http.HandleFunc("/users/", func(w http.ResponseWriter, r *http.Request) {
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
	})

	fmt.Println("Сервер запущен на http://localhost:8080")
	fmt.Println("GET /users - получить всех пользователей")
	fmt.Println("GET /users/1 - получить пользователя с ID=1")
	fmt.Println("POST /users - создать пользователя")
	fmt.Println("DELETE /users/1 - удалить пользователя")
	fmt.Println()
	fmt.Println("Пример POST запроса:")
	fmt.Println(`curl -X POST http://localhost:8080/users/ -H "Content-Type: application/json" -d '{"name":"Anna","email":"anna@mail.com","age":28}'`)

	http.ListenAndServe(":8080", nil)
}
