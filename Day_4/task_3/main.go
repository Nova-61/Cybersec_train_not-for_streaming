package main

import "fmt"

func calculate(a, b float64, op string) (result float64, err error) {

	if b == 0 && op == "/" {
		result = 0
		err = fmt.Errorf("division by zero")
		return
	} else {
		if op == "+" {
			result = a + b
			err = nil
		} else if op == "-" {
			result = a - b
			err = nil
		} else if op == "*" {
			result = a * b
			err = nil
		} else if op == "/" {
			result = a / b
			err = nil
		} else {
			result = 0
			err = fmt.Errorf("unsupported operation: %s", op)
		}
	}
	return
}

func main() {
	fmt.Println(calculate(4, 2, "+"))
}
