package main

import (
	"fmt"
)

func Add(a, b int) int {
	return a + b
}

func Multiply(a, b int) int {
	return a * b
}

func Divide(a, b int) (int, error) {
	switch b {
	case 0:
		return 0, fmt.Errorf("деление на ноль")
	default:
		return a / b, nil
	}
}

func main() {
	result1 := Add(2, 3)
	fmt.Printf("Add: %d\n", result1)

	result2 := Multiply(2, 3)
	fmt.Printf("Multiply: %d\n", result2)

	result3, err := Divide(2, 3)
	if err != nil {
		fmt.Println("Ошибка:", err)
	} else {
		fmt.Printf("Divide: %d\n", result3)
	}

	// Проверка деления на ноль
	result4, err := Divide(10, 0)
	if err != nil {
		fmt.Println("Ошибка:", err)
	} else {
		fmt.Printf("Divide: %d\n", result4)
	}
}
