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

func worker(id int, tasks <-chan Task, results chan<- Result) {
	defer wg.Done()
	for task := range tasks {
		result := Result{
			TaskID: task.ID,
			Square: task.Value * task.Value,
		}
		results <- result
	}
}

func main() {
	tasks := make(chan Task, 20)
	results := make(chan Result, 20)

	// Запускаем 5 воркера
	for i := 1; i <= 5; i++ {
		wg.Add(1)
		go worker(i, tasks, results)
	}

	// Отправляем 17 задач
	for i := 1; i <= 17; i++ {
		tasks <- Task{
			ID:    i,
			Value: i,
		}
	}
	close(tasks)

	// Закрываем результаты после завершения всех воркеров
	go func() {
		wg.Wait()
		close(results)
	}()

	// Собираем и выводим результаты
	for result := range results {
		fmt.Printf("Task %d: %d^2 = %d\n", result.TaskID, result.TaskID, result.Square)
	}
}
