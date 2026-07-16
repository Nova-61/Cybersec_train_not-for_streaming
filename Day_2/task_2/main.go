package main

import "fmt"

func main() {
	const TaxRate = 0.20
	const (
		Admin = iota
		Manager
		User
		Guest
	)
	fmt.Printf(" %f, %d, %d, %d, %d", TaxRate, Admin, Manager, User, Guest)
}
