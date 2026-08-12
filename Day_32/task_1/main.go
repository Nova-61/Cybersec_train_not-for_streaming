package main

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"
)

func main() {
	// Создаём цель (куда перенаправлять)
	target, _ := url.Parse("http://localhost:8080")

	// Создаём прокси
	proxy := httputil.NewSingleHostReverseProxy(target)

	client := &http.Client{Timeout: 2 * time.Second}
	resp, err := client.Get("http://localhostL8080/")
	if err != nil {
		fmt.Printf("Ошибка: %s\n", err)
	} else {
		defer resp.Body.Close()
		fmt.Println("Сервер успшено работает")
	}
	// Обработчик для прокси
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("[PROXY] %s %s\n", r.Method, r.URL.Path)
		proxy.ServeHTTP(w, r)
	})

	// Запускаем сервер
	fmt.Println("Прокси запущен на :8888")
	fmt.Println("Перенаправляет на :8080")
	http.ListenAndServe(":8888", proxy)
}
