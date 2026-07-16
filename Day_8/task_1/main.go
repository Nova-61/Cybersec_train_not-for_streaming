package main

import (
	"fmt"
	"time"
)

func printNumbers() {
	for i := 1; i <= 5; i++ {
		fmt.Println(i)
	}
}

func main() {
	go printNumbers()

	time.Sleep(1 * time.Second)
	fmt.Println("Hello go runtine")

}
