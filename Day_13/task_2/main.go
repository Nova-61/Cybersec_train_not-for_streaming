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
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	result1 := slowOperation(ctx, 1*time.Second)
	fmt.Println("Результат 1 (1 сек):", result1)

	result2 := slowOperation(ctx, 3*time.Second)
	fmt.Println("Результат 2 (3 сек):", result2)
}
