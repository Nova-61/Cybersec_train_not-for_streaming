package main

import "fmt"

type Mover interface {
	Move() string
}

type Car struct {
	Name string
}

func (c Car) Move() string {
	return fmt.Sprintf("%s is moving on wheels", c.Name)
}

type Bicycle struct {
	Name string
}

func (b Bicycle) Move() string {
	return fmt.Sprintf("%s is moving on pedals", b.Name)
}

func main() {
	car := Car{Name: "Toyota"}
	bicycle := Bicycle{Name: "Giant"}

	fmt.Println(car.Move())
	fmt.Println(bicycle.Move())
}
