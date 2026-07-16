package main

import (
	"fmt"
	"sync"
)

type Operation struct {
	ID     int
	A      float64
	B      float64
	Op     string
	Result float64
	Error  error
}

func calculate(op *Operation, wg *sync.WaitGroup, mu *sync.Mutex) {
	defer wg.Done()

	var result float64
	var err error

	switch op.Op {
	case "+":
		result = op.A + op.B
	case "-":
		result = op.A - op.B
	case "*":
		result = op.A * op.B
	case "/":
		if op.B == 0 {
			err = fmt.Errorf("деление на ноль")
		} else {
			result = op.A / op.B
		}
	default:
		err = fmt.Errorf("неизвестная операция: %s", op.Op)
	}

	mu.Lock()
	op.Result = result
	op.Error = err
	mu.Unlock()
}

func main() {
	operations := []Operation{
		{ID: 1, A: 10, B: 5, Op: "+"},
		{ID: 2, A: 8, B: 3, Op: "-"},
		{ID: 3, A: 7, B: 2, Op: "*"},
		{ID: 4, A: 10, B: 0, Op: "/"},
		{ID: 5, A: 15, B: 3, Op: "/"},
		{ID: 6, A: 6, B: 2, Op: "+"},
		{ID: 7, A: 20, B: 4, Op: "/"},
		{ID: 8, A: 9, B: 3, Op: "-"},
		{ID: 9, A: 5, B: 5, Op: "*"},
		{ID: 10, A: 100, B: 0, Op: "/"},
	}

	var wg sync.WaitGroup
	var mu sync.Mutex

	for i := range operations {
		wg.Add(1)
		go calculate(&operations[i], &wg, &mu)
	}

	wg.Wait()

	for _, op := range operations {
		if op.Error != nil {
			fmt.Printf("%.0f %s %.0f = Ошибка: %v\n", op.A, op.Op, op.B, op.Error)
		} else {
			fmt.Printf("%.0f %s %.0f = %.0f\n", op.A, op.Op, op.B, op.Result)
		}
	}
}
