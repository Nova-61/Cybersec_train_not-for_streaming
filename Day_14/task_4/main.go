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

func worker(id int, tasks <-chan Task) <-chan Result {
	out := make(chan Result)

	go func() {
		defer close(out)
		for task := range tasks {
			result := Result{
				TaskID: task.ID,
				Square: task.Value * task.Value,
			}
			out <- result
		}
	}()

	return out
}

func generateTask(n int) <-chan Task {
	out := make(chan Task)

	go func() {
		defer close(out)
		for i := 1; i <= n; i++ {
			out <- Task{
				ID:    i,
				Value: i,
			}
		}
	}()

	return out
}

func main() {
	// Генерируем задачи
	tasks := generateTask(10)

	// Запускаем 3 воркера (каждый получает свой канал)
	worker1 := worker(1, tasks)
	worker2 := worker(2, tasks)
	worker3 := worker(3, tasks)

	// Объединяем результаты в один канал
	results := fanIn(worker1, worker2, worker3)

	// Выводим результаты
	for result := range results {
		fmt.Printf("Task %d: %d^2 = %d\n", result.TaskID, result.TaskID, result.Square)
	}
}
