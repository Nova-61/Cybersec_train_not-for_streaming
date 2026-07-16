package main

import (
	"fmt"
)

func generate(n int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for i := 1; i <= n; i++ {
			out <- i
		}
	}()
	return out
}

func multiply(in <-chan int, factor int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for n := range in {
			out <- n * factor
		}
	}()
	return out
}

func filterEven(in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for n := range in {
			if n%2 == 0 {
				out <- n
			}
		}
	}()
	return out
}

func sum(in <-chan int) int {
	total := 0
	for n := range in {
		total += n
	}
	return total
}

func main() {
	numbers := generate(10)
	multiplied := multiply(numbers, 2)
	filtered := filterEven(multiplied)

	result := sum(filtered)
	fmt.Println("Sum of even numbers multiplied by 2:", result)

}
