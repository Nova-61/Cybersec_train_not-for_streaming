package main

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

// Proxy - структура для хранения целевого сервера
type Proxy struct {
	target *url.URL
}

// NewProxy - создаёт новый прокси
func NewProxy(targetURL string) (*Proxy, error) {
	target, err := url.Parse(targetURL)
	if err != nil {
		return nil, err
	}
	return &Proxy{target: target}, nil
}

// ServeHTTP - обрабатывает HTTP-запросы
func (p *Proxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Логируем запрос
	fmt.Printf("[PROXY] %s %s\n", r.Method, r.URL.Path)

	// Проверяем заголовки (безопасность)
	if strings.Contains(r.Header.Get("Authorization"), "secret") {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	// Создаём reverse proxy
	proxy := httputil.NewSingleHostReverseProxy(p.target)

	// Передаём запрос
	proxy.ServeHTTP(w, r)
}

func main() {
	proxy, err := NewProxy("http://localhost:8080")
	if err != nil {
		fmt.Println("Ошибка:", err)
		return
	}

	fmt.Println("Прокси запущен на :8888")
	fmt.Println("Перенаправляет на http://localhost:8080")

	err = http.ListenAndServe(":8888", proxy)
	if err != nil {
		fmt.Println("Ошибка запуска сервера:", err)
	}
}

// Браузер/Клиент
//      |
//      |  http://localhost:8888/users
//      ↓
// ┌─────────────────────────────┐
// │   ПРОКСИ (порт 8888)        │
// │                              │
// │  1. Логирует запрос          │
// │  2. Проверяет заголовки      │
// │  3. Меняет Host на target    │
// └─────────────────────────────┘
//      |
//      |  http://localhost:8080/users
//      ↓
// ┌─────────────────────────────┐
// │   API СЕРВЕР (порт 8080)    │
// │                              │
// │  Обрабатывает запрос         │
// │  Возвращает ответ            │
// └─────────────────────────────┘
//      |
//      |  Ответ
//      ↓
//      Клиент
