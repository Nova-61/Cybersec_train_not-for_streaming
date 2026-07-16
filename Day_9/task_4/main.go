package main

import (
	"fmt"
	"sync"
)

func producer(ch chan int, wg *sync.WaitGroup) {
	defer wg.Done()
	for i := 0; i < 10; i++ {
		ch <- i
	}
	close(ch)
}

func main() {
	ch := make(chan int)
	var wg sync.WaitGroup
	wg.Add(1)
	go producer(ch, &wg)

	wg.Wait()

	for value := range ch {
		fmt.Println("Received:", value)
	}
}
