package main

import (
	"fmt"
	"math/rand"
	"time"
)

func main() {
	rand.Seed(time.Now().UnixNano())
	
	secretNumber := rand.Intn(100) + 1
	var guess int
	attempts := 0

	fmt.Println("Угадайте число от 1 до 100")
	
	for {
		fmt.Print("Введите ваше предположение: ")
		fmt.Scan(&guess)
		attempts++

		if guess < secretNumber {
			fmt.Println("Слишком маленькое число. Попробуйте снова.")
		} else if guess > secretNumber {
			fmt.Println("Слишком большое число. Попробуйте снова.")
		} else {
			fmt.Printf("Поздравляем! Вы угадали число %d за %d попыток.\n", secretNumber, attempts)
			break
		}
	}
}