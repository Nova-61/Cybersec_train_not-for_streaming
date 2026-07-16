package main

import (
	"fmt"
	"net/http"
)

func myMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("1. Запрос получен!")
		next(w, r)
		fmt.Println("3. Запрос обработан!")
	}
}

func loggerMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Запрос начат")
		next(w, r)
		fmt.Println("Запрос завершен")
	}
}


func helloHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("2. Привет, мир!")
	fmt.Fprintf(w, "Hello, World!")
}

func pingHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Ping handler")
	fmt.Fprintf(w, "pong")
}

func main() {
	http.HandleFunc("/", myMiddleware(helloHandler))
	http.HandleFunc("/ping", loggerMiddleware(pingHandler))

	fmt.Println("Server on http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}