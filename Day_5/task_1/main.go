package main

import "fmt"

type Book struct {
	Title  string
	author string
	Pages  int
	Year   int
}

func main() {

	fmt.Println(Book{Title: "Kill the bill"})
	fmt.Println(Book{Title: "The Great Gatsby", author: "F. Scott Fitzgerald", Pages: 180, Year: 1925})
}
