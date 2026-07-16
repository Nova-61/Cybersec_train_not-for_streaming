package main

import (
	"fmt"
	"time"
)

func printNumbers() {
	for i := 1; i <= 1000000; i++ {
		fmt.Println(i)
	}
}

func main() {
	go printNumbers()

	time.Sleep(2 * time.Second)
	fmt.Println("Hello go runtine")

}
