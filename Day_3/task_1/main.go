package main

import "fmt"

func main() {
	var a int
	fmt.Scan(&a)
	if a < 18 {
		fmt.Println("Доступ запрещен")
	} else if a >= 18 && a < 65 {
		fmt.Println("Доступ разрешен")
	} else if a >= 65 {
		fmt.Println("Пенсионер, доступ с ограмничениями")
	}
}
