package main

import (
	"context"
	"fmt"
	"sync"
)

var wg sync.WaitGroup

func worker(ctx context.Context) {
	defer wg.Done()

	for{
		select{
			case <- ctx.Done():
				fmt.Println("Worker: остановлен по сигналу отмены")
			default:
				fmt.Println("Работаю...")
			}
	}
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	wg.Add(1)
	go worker(ctx)

	cancel()

	wg.Wait()
}