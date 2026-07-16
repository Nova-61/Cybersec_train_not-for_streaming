package main

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

// Функция которая долго обрабатывает данные (5 секунд)
func processData(ctx context.Context, data string) string {
	// Проверяем, не отменили ли нас
	select {
	case <-ctx.Done():
		return "Отменено"
	default:
		return "Обработано: " + data
	}
}

// Обработчик запроса
func handler(w http.ResponseWriter, r *http.Request) {
	// Получаем параметр data из URL
	data := r.URL.Query().Get("data")
	if data == "" {
		http.Error(w, "Нет параметра data", http.StatusBadRequest)
		return
	}

	// Создаём контекст с таймаутом 3 секунды
	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	// Запускаем обработку в отдельной горутине
	resultCh := make(chan string)
	go func() {
		resultCh <- processData(ctx, data)
	}()

	// Ждём результат или таймаут
	select {
	case result := <-resultCh:
		if result == "Отменено" {
			http.Error(w, "Время вышло", 408)
		} else {
			fmt.Fprintf(w, "Результат: %s", result)
		}
	case <-ctx.Done():
		http.Error(w, "Время вышло", 408)
	}
}

func main() {
	// Запускаем сервер
	http.HandleFunc("/process", handler)
	fmt.Println("Сервер на http://localhost:8080")
	fmt.Println("Проверка: http://localhost:8080/process?data=hello")
	http.ListenAndServe(":8080", nil)
}
