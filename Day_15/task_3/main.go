package main

import (
	"fmt"
	"net/http"
)

// Простые структуры
type User struct {
	ID   int
	Name string
	Age  int
}

type Product struct {
	ID    int
	Name  string
	Price float64
}

// Список юзеров
var users = []User{
	{ID: 1, Name: "Алексей", Age: 25},
	{ID: 2, Name: "Мария", Age: 30},
	{ID: 3, Name: "Иван", Age: 22},
}

// Список товаров
var products = []Product{
	{ID: 1, Name: "Ноутбук", Price: 999.99},
	{ID: 2, Name: "Мышь", Price: 29.99},
	{ID: 3, Name: "Клавиатура", Price: 49.99},
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Главная страница")
}

func UsersHandler(w http.ResponseWriter, r *http.Request) {
	// Просто выводим список юзеров
	fmt.Fprintf(w, "Список пользователей:\n")
	for _, u := range users {
		fmt.Fprintf(w, "ID: %d, Имя: %s, Возраст: %d\n", u.ID, u.Name, u.Age)
	}
}

func ProductsHandler(w http.ResponseWriter, r *http.Request) {
	// Просто выводим список товаров
	fmt.Fprintf(w, "Список товаров:\n")
	for _, p := range products {
		fmt.Fprintf(w, "ID: %d, Название: %s, Цена: %.2f\n", p.ID, p.Name, p.Price)
	}
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("Запрос: %s %s\n", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/", HomeHandler)
	mux.HandleFunc("/users", UsersHandler)
	mux.HandleFunc("/products", ProductsHandler)

	handler := loggingMiddleware(mux)

	fmt.Println("Server on http://localhost:8080")
	http.ListenAndServe(":8080", handler)
}
