package main

import "fmt"

func calculate(a, b float64, op string) (float64, error) {

	if (a == 0 || b == 0) && op == "/" {
		return 0, fmt.Errorf("division by zero")
	} else {
		if op == "+" {
			return a + b, nil
		} else if op == "-" {
			return a - b, nil
		} else if op == "*" {
			return a * b, nil
		} else if op == "/" {
			return a / b, nil
		} else {
			return 0, fmt.Errorf("unsupported operation: %s", op)
		}
	}
}

func main() {
	fmt.Println(calculate(4, 2, "+"))
}
