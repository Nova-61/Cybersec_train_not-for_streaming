package main

import (
	"fmt"
	"math"
)

type Shape interface {
	Area() float64
	Perimeter() float64
}

type Circle struct {
	Radius float64
}

type Square struct {
	Side float64
}

type Stringer interface {
	String() string
}

func (c Circle) Area() float64 {
	return math.Pi * c.Radius * c.Radius
}

func (c Circle) Perimeter() float64 {
	return 2 * math.Pi * c.Radius
}

func (c Circle) String() string {
	return fmt.Sprintf("Circle with radius: %.2f", c.Radius)
}

func (s Square) Area() float64 {
	return s.Side * s.Side
}

func (s Square) Perimeter() float64 {
	return 4 * s.Side
}

func (s Square) String() string {
	return fmt.Sprintf("Square with side: %.2f", s.Side)
}

func PrintAnything(v interface{}) {
	switch val := v.(type) {
	case int:
		fmt.Printf("Число: %d\n", val)
	case string:
		fmt.Printf("Строка: %s\n", val)
	case bool:
		fmt.Printf("Булево: %t\n", val)
	default:
		fmt.Printf("Неизвестный тип: %T\n", val)
	}
}

func main() {
	shapes := []Shape{
		Circle{Radius: 5},
		Square{Side: 4},
	}

	for _, shape := range shapes {
		fmt.Printf("Area: %.2f, Perimeter: %.2f\n", shape.Area(), shape.Perimeter())
	}

	fmt.Println("\n--- Проверка PrintAnything ---")

	PrintAnything(42)        // Число: 42
	PrintAnything("hello")   // Строка: hello
	PrintAnything(true)      // Булево: true
	PrintAnything(3.14)      // Неизвестный тип: float64
	PrintAnything(Circle{5}) // Неизвестный тип: main.Circle
	PrintAnything(shapes)    // Неизвестный тип: []main.Shape

	fmt.Println("\n--- Проверка Stringer ---")
	var s Stringer = Circle{Radius: 5}
	fmt.Println(s.String()) // Circle with radius: 5.00

	s = Square{Side: 4}
	fmt.Println(s.String()) // Square with side: 4.00
}
