package main

import "fmt"

func addBonus(bonus float64, account *float64) {
	if account == nil {
		fmt.Println("Ошибка с указателем")
		return
	}
	*account += bonus
}

func main() {
	var balance float64 = 1000.0
	var ptr *float64 = &balance

	fmt.Printf("Баланс до бонуса: %.2f\n", balance)

	var add float64
	fmt.Scan(&add)

	addBonus(add, ptr)

	fmt.Printf("Баланс до бонуса: %.2f\n", balance)

	fmt.Println(balance)
}
