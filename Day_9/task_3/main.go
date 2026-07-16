package main

import "fmt"

func gen(nums ...int) <-chan int {
	out := make(chan int)
	go func() {
		for _, n := range nums {
			out <- n
		}
		close(out)
	}()
	return out
}

func square(in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		for n := range in {
			out <- n * n
		}
		close(out)
	}()
	return out
}

func print(in <-chan int) {
	for n := range in {
		fmt.Println(n)
	}
}

func main() {
	numbers := gen(1, 2, 3, 4, 5)
	squares := square(numbers)
	print(squares)
}
