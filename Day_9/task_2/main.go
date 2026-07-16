package main

import (
	"fmt"
	"sync"
)

func main() {
	ch := make(chan int, 3)
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		ch <- 42
		ch <- 43
		ch <- 44
		ch <- 45
	}()

	value := <-ch
	fmt.Println("Received:", value)

	wg.Wait()
	close(ch)

	fmt.Println("Остальные значнения из канала:")
	for i := range ch {
		fmt.Println("Received:", i)
	}

}
