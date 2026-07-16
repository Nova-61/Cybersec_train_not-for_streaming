package main

import (
	"fmt"
	"runtime"
)

func main() {
	fmt.Println("Hello, World!")
	fmt.Printf("Go version: %s\n", runtime.Version())
	fmt.Printf("Operating System: %s\n", runtime.GOOS)
	fmt.Printf("Architecture: %s\n", runtime.GOARCH)
}


