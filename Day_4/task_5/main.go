package main

import "fmt"

var history []string

func calc(a, b float64, op string) (float64, error) {
	if b == 0 && op == "/" {
		return 0, fmt.Errorf("division by zero")
	}

	switch op {
	case "+":
		return a + b, nil
	case "-":
		return a - b, nil
	case "*":
		return a * b, nil
	case "/":
		return a / b, nil
	default:
		return 0, fmt.Errorf("unknown operation")
	}
}

func createCalculator() func(a, b float64, op string) (float64, error) {
	history = []string{}

	return func(a, b float64, op string) (float64, error) {
		result, err := calc(a, b, op)

		if err == nil {
			history = append(history, fmt.Sprintf("%.2f %s %.2f = %.2f", a, op, b, result))
		} else {
			history = append(history, fmt.Sprintf("%.2f %s %.2f = Ошибка: %v", a, op, b, err))
		}

		return result, err
	}
}

func getHistory() []string {
	return history
}

func main() {
	calcWithHistory := createCalculator()

	calcWithHistory(3, 2, "+")
	calcWithHistory(5, 7, "*")
	calcWithHistory(10, 2, "/")
	calcWithHistory(15, 4, "-")
	calcWithHistory(8, 0, "/")

	fmt.Println("История операций:")
	for i, op := range getHistory() {
		fmt.Printf("%d. %s\n", i+1, op)
	}
}
