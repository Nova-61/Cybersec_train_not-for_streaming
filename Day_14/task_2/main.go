package main

import (
	"fmt"
)

type Task struct {
	ID    int
	Value int
}

func generateTask(n int) <-chan Task {
	out := make(chan Task)

	go func() {
		for i := 1; i <= n; i++ {
			out <- Task{
				ID:    i,
				Value: i,
			}
		}
		close(out)
	}()
	return out
}

func main() {
	tasks := generateTask(5)

	for task := range tasks {
		fmt.Printf("Task %d: %d\n", task.ID, task.Value)
	}
}
