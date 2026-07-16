package main

import "fmt"

const (
	USDToRUB = 90.0
	EURToRUB = 100.0
	RUBToUSD = 1.0 / USDToRUB // 0.011...
	RUBToEUR = 1.0 / EURToRUB
)

// Конвертирует сумму из одной валюты в другую
func convert(amount float64, from string, to string) float64 {
	// Сначала переводим всё в рубли
	var rubles float64

	switch from {
	case "USD":
		rubles = amount * USDToRUB
	case "EUR":
		rubles = amount * EURToRUB
	case "RUB":
		rubles = amount
	default:
		fmt.Println("Неизвестная валюта:", from)
		return 0
	}

	// Теперь из рублей в нужную валюту
	switch to {
	case "USD":
		return rubles * RUBToUSD
	case "EUR":
		return rubles * RUBToEUR
	case "RUB":
		return rubles
	default:
		fmt.Println("Неизвестная валюта:", to)
		return 0
	}
}

// Функция с указателем
func convertWithPointer(amount *float64, from string, to string) {
	if amount == nil {
		fmt.Println("Ошибка: указатель nil")
		return
	}
	*amount = convert(*amount, from, to)
}

func main() {
	// Без указателя
	result := convert(100, "USD", "RUB")
	fmt.Printf("100 USD = %.2f RUB\n", result)

	// С указателем
	balance := 50.0
	fmt.Printf("До конвертации: %.2f EUR\n", balance)
	convertWithPointer(&balance, "EUR", "USD")
	fmt.Printf("После конвертации: %.2f USD\n", balance)

	// Проверка на ошибки
	bad := convert(100, "XXX", "RUB")
	fmt.Printf("Результат с ошибкой: %.2f\n", bad)
}
