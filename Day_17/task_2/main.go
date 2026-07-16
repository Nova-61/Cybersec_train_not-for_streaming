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
	jsonString := `{"id":2,"name":"Phone","price":599.99,"in_stock":false}`
	var product Product
	err := json.Unmarshal([]byte(jsonString), &product)
	if err != nil {
		fmt.Println("Ошибка:", err)
		return
	}
	fmt.Println(product.Name)
}