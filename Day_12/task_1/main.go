package main

import (
	"fmt"
	"sync"
	"time"
)

func worker(id int, jobs int) {
	fmt.Printf("Worker %d: started\n", id)
	time.Sleep(time.Duration(jobs) * time.Second)
	fmt.Printf("Worker %d: finished\n", id)
}

func main() {
	var wg sync.WaitGroup

	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(n int) {
			defer wg.Done()
			worker(n, n+1)
		}(i)
	}

	wg.Wait()
	fmt.Println("Все горутины завершены")
}