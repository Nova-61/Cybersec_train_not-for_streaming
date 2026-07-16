package main

import (
	"fmt"
	"net/http"
	"strconv"
)

func handler_hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Привет Мир!")
}

func handler_with_name(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	if name == "" {
		name = "Гость"
	}
	fmt.Fprintf(w, "Привет %s!", name)
}

func handler_calc(w http.ResponseWriter, r *http.Request) {
	aStr := r.URL.Query().Get("a")
	bStr := r.URL.Query().Get("b")
	op := r.URL.Query().Get("op")

	if aStr == "" || bStr == "" || op == "" {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"error": "missing parameters"}`)
		return
	}

	a, errA := strconv.Atoi(aStr)
	b, errB := strconv.Atoi(bStr)

	if errA != nil || errB != nil {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"error": "invalid numbers"}`)
		return
	}

	var result int

	switch op {
	case "+":
		result = a + b
	case "-":
		result = a - b
	case "*":
		result = a * b
	case "/":
		if b == 0 {
			w.Header().Set("Content-Type", "application/json")
			fmt.Fprintf(w, `{"error": "division by zero"}`)
			return
		}
		result = a / b
	default:
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"error": "unknown operation"}`)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, `{"result": %d}`, result)
}

func main() {
	http.HandleFunc("/hello", handler_hello)
	http.HandleFunc("/greet", handler_with_name)
	http.HandleFunc("/calc", handler_calc)

	fmt.Println("Server on http://localhost:8080")
	fmt.Println("  /hello               - Привет Мир!")
	fmt.Println("  /greet?name=Ivan     - Привет Ivan!")
	fmt.Println("  /calc?a=10&b=5&op=+  - 15")
	http.ListenAndServe(":8080", nil)
}