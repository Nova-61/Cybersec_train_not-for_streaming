package main

import (
	"encoding/json"
	"fmt"
)

type Product struct {
	ID       int
	Name     string
	Price    float64
	InStock  bool
}

func main() {
	product := Product{
		ID:       1,
		Name:     "Laptop",
		Price:    999.99,
		InStock:  true,
	}

	jsonData, err := json.Marshal(product)
	if err != nil {
		fmt.Println("Ошибка:", err)
		return
	}

	fmt.Println(string(jsonData))
}