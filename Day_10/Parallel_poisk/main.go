package main

import (
	"fmt"
	"sync"
)

func isPrime(n int) bool {
	if n < 2 {
		return false
	}
	for i := 2; i*i <= n; i++ {
		if n%i == 0 {
			return false
		}
	}
	return true
}

func main() {
	numbers := []int{2, 3, 4, 5, 6, 7, 8, 9, 10}
	var wg sync.WaitGroup

	for _, n := range numbers {
		wg.Add(1)
		go func(num int) {
			defer wg.Done()
			fmt.Printf("%d: %t\n", num, isPrime(num))
		}(n)
	}

	wg.Wait()
}
