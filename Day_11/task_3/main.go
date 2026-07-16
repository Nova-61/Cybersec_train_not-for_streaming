package main

import (
	"fmt"
	"sync"
)

func RCounter(wg *sync.WaitGroup, mu *sync.RWMutex, data map[string]int) {
	defer wg.Done()
	mu.RLock()
	_ = data["counter"]
	mu.RUnlock()
}

func WCounter(wg *sync.WaitGroup, mu *sync.RWMutex, data map[string]int) {
	defer wg.Done()
	mu.Lock()
	data["counter"]++
	mu.Unlock()
}

func main() {
	data := make(map[string]int)
	data["counter"] = 0

	var rw sync.RWMutex
	var wg sync.WaitGroup

// Читатель
	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go RCounter(&wg, &rw, data)
	}
// Писатель
	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go WCounter(&wg, &rw, data)
	}

	wg.Wait()
	fmt.Println(data["counter"])
}
