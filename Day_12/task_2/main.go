package main

import (
	"fmt"
	"sync"
)

var config map[string]string

func loadConfig() {
	fmt.Println("Loading config...")
	config = make(map[string]string)
	config["host"] = "localhost"
	config["port"] = "8080"
	config["user"] = "admin"
}

func main() {
	var once sync.Once

	for i := 0; i < 10; i++ {
		go func() {
			once.Do(loadConfig)
		}()
	}

	fmt.Println("Config:", config)
}
