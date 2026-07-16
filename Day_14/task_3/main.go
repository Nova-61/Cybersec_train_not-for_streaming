package main

import (
	"fmt"
	"sync"
)

var wg sync.WaitGroup

type Task struct {
	ID    int
	Value int
}

type Result struct {
	TaskID int
	Square int
}

func fanIn(results ...<-chan Result) <-chan Result {
	out := make(chan Result)

	for _, ch := range results {
		wg.Add(1)
		go func(c <-chan Result) {
			defer wg.Done()
			for result := range c {
				out <- result
			}
		}(ch)
	}

	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}

func main() {

	// Объединяем результаты в один канал
	results := fanIn()

	// Выводим результаты
	for result := range results {
		fmt.Printf("Task %d: %d^2 = %d\n", result.TaskID, result.TaskID, result.Square)
	}
}
