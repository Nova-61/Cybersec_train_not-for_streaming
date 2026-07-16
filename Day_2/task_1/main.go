package main

import "fmt"

func main() {
	var name string = "Ivan"
	var age = 30
	city := "Komsa"
	var my_salary float64 = 1234.1234
	var isActive bool = true
	
	fmt.Printf("Hi, my name is %s. I am %d years old. I live in %s. My salary is %.2f. Am I active? %t\n", 
		name, age, city, my_salary, isActive)
}
