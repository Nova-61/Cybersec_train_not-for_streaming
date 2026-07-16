package main

import (
	"fmt"
	"sync"
)

func Counter(wg *sync.WaitGroup, mu *sync.Mutex, counter *int) {
	defer wg.Done()
	mu.Lock()
	*counter++
	mu.Unlock()
}

func main() {
	counter := 0
	var wg sync.WaitGroup
	var mu sync.Mutex

	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go Counter(&wg, &mu, &counter)
	}
	wg.Wait()
	fmt.Println(counter)
}
