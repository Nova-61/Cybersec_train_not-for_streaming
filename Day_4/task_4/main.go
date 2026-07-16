package main

import "fmt"

func multyplier(factor int) func(int) int {
	return func(x int) int {
		return factor * x
	}
}

func main() {
	fmt.Println(multyplier(2)(3))
}
