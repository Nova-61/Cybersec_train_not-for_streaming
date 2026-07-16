package main

import (
	"fmt"
	"sync"
)

func main() {
	ch := make(chan int)
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		ch <- 42
	}()

	value := <-ch
	fmt.Println("Received:", value)

	wg.Wait()
	close(ch)
}
