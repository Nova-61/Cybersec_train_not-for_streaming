package main

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

func loginHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		fmt.Fprintf(w, `<form method="POST">
			<input name="username" placeholder="Логин">
			<input name="password" type="password" placeholder="Пароль">
			<button type="submit">Войти</button>
		</form>`)

	case "POST":
		if err := r.ParseForm(); err != nil {
			http.Error(w, "Ошибка обработки формы", http.StatusBadRequest)
			return
		}

		username := r.Form.Get("username")
		password := r.Form.Get("password")

		if username == "" || password == "" {
			http.Error(w, "Все поля обязательны", http.StatusBadRequest)
			return
		}

		if username == "admin" && password == "secret" {
			fmt.Fprintf(w, "Успешный вход!")
		} else {
			http.Error(w, "Неверный логин или пароль", http.StatusUnauthorized)
		}

	default:
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
	}
}

func main() {
	http.HandleFunc("/login", loginHandler)

	srv := &http.Server{
		Addr:    ":8080",
		Handler: nil,
	}

	go func() {
		fmt.Println("Server running on http://localhost:8080")
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

	fmt.Println("Server successfully stopped")
}
