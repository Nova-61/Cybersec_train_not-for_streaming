package main

import (
	"fmt"
	"net/http"
	"strconv"
)

func greetHandler(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	if name == "" {
		name = "Гость"
	}
	fmt.Fprintf(w, "Привет, %s!", name)
}

func addHandler(w http.ResponseWriter, r *http.Request) {
	aStr := r.URL.Query().Get("a")
	bStr := r.URL.Query().Get("b")

	if aStr == "" || bStr == "" {
		http.Error(w, "Параметры a и b обязательны", http.StatusBadRequest)
		return
	}

	a, errA := strconv.Atoi(aStr)
	b, errB := strconv.Atoi(bStr)

	if errA != nil || errB != nil {
		http.Error(w, "Параметры должны быть числами", http.StatusBadRequest)
		return
	}

	result := a + b
	fmt.Fprintf(w, "%d + %d = %d", a, b, result)
}

func userHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")

	if id == "" {
		fmt.Fprintf(w, "Все пользователи")
		return
	}

	fmt.Fprintf(w, "Пользователь с ID: %s", id)
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("Запрос: %s %s\n", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/greet", greetHandler)
	mux.HandleFunc("/add", addHandler)
	mux.HandleFunc("/user", userHandler)

	handler := loggingMiddleware(mux)

	fmt.Println("Server running on http://localhost:8080")
	fmt.Println("  /greet?name=Ivan  - Приветствие")
	fmt.Println("  /add?a=5&b=3      - Сложение")
	fmt.Println("  /user?id=1        - Информация о пользователе")

	http.ListenAndServe(":8080", handler)
}