package main

import (
	"fmt"
	"sync"
)

type Counter interface {
	Increment()
	Decrement()
	Value() int
}

type SafeCounter struct {
	mu    sync.Mutex
	value int
}

func (s *SafeCounter) Increment() {
	s.mu.Lock()
	s.value++
	s.mu.Unlock()
}

func (s *SafeCounter) Decrement() {
	s.mu.Lock()
	s.value--
	s.mu.Unlock()
}

func (s *SafeCounter) Value() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.value
}

func RunOperations(c Counter, increment, decrement int, wg *sync.WaitGroup) {
	defer wg.Done()
	for i := 0; i < increment; i++ {
		c.Increment()
	}
	for i := 0; i < decrement; i++ {
		c.Decrement()
	}
}

func main() {
	var counter Counter = &SafeCounter{}
	var wg sync.WaitGroup

	for i := 0; i < 50; i++ {
		wg.Add(1)
		go RunOperations(counter, 10, 0, &wg)
	}

	wg.Wait()
	fmt.Println(counter.Value())
}
