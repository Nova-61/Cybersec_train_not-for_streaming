package main

import (
	"fmt"
	"time"
)

func getWithTimeout(timeout time.Duration) (string, error) {
	result := make(chan string)

	go func() {
		time.Sleep(2 * time.Second)
		result <- "готово"
	}()

	select {
	case res := <-result:
		return res, nil
	case <-time.After(timeout):
		return "", fmt.Errorf("timeout after %v", timeout)
	}
}

func main() {
	for _, timeout := range []time.Duration{time.Second, 3 * time.Second} {
		result, err := getWithTimeout(timeout)
		if err != nil {
			fmt.Printf("timeout %v: %v\n", timeout, err)
			continue
		}

		fmt.Printf("result with timeout %v: %s\n", timeout, result)
	}
}