package main

import (
	"fmt"
	"sync"
	"time"
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

type PrimeResult struct {
	Number   int
	IsPrime  bool
	Duration time.Duration
	WorkerID int
}

func worker(id int, numbers <-chan int, results chan<- PrimeResult, wg *sync.WaitGroup) {
	defer wg.Done()

	for n := range numbers {
		start := time.Now()
		isPrimeResult := isPrime(n)
		duration := time.Since(start)
		
		results <- PrimeResult{
			Number:   n,
			IsPrime:  isPrimeResult,
			Duration: duration,
			WorkerID: id,
		}
	}
}

func checkNumbers(numbers []int) map[int]bool {
	result := make(map[int]bool)

	var mu sync.Mutex
	var wg sync.WaitGroup

	for _, n := range numbers {
		wg.Add(1)
		go func(num int) {
			defer wg.Done()
			mu.Lock()
			result[num] = isPrime(num)
			mu.Unlock()
		}(n)
	}
	wg.Wait()
	return result
}

func main() {
	numbers := []int{2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20}
	result := checkNumbers(numbers)
	for n, isPrime := range result {
		fmt.Printf("Number %d is prime: %t\n", n, isPrime)
	}
}
