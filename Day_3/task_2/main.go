package main

import "fmt"

func main() {
	month := 3

	switch month {
	case 12, 1, 2:
		fmt.Println("Зима")
	case 3, 4, 5:
		fmt.Println("Весна")
	case 6, 7, 8:
		fmt.Println("Лето")
	case 9, 10, 11:
		fmt.Println("Осень")
	default:
		fmt.Println("Неверный номер месяца")
	}
}
