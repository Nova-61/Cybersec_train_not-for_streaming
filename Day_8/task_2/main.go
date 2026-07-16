package main

import "fmt"

func main() {
	for num := 0; num < 5; num++ {
		go func(n int) {
			fmt.Printf("Горутина %d: работает!\n", n)
		}(num)
	}
}
