package main

import (
	"fmt"
)

func main() {
	ch := make(chan int)

	select {
	case msg := <-ch:
		fmt.Println("Received from channel:", msg)
	default:
		fmt.Println("No message received from channel")
	}
}
