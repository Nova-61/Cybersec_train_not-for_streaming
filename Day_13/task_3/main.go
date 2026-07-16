package main

import (
	"context"
	"fmt"
	"time"
)

func slowOperation(ctx context.Context, duration time.Duration) string {
	select {
	case <-time.After(duration):
		return "Операция выполнена"
	case <-ctx.Done():
		return "Операция отменена"
	}
}

func main() {
	// ✅ Устанавливаем дедлайн на 3 секунды от текущего времени
	deadline := time.Now().Add(3 * time.Second)
	ctx, cancel := context.WithDeadline(context.Background(), deadline)
	defer cancel()

	// ✅ Запускаем slowOperation с длительностью 5 секунд
	result := slowOperation(ctx, 5*time.Second)
	fmt.Println("Результат:", result)

	// ✅ Выводим причину отмены
	fmt.Println("Причина отмены:", ctx.Err())
}
