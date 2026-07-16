package main

import "fmt"

type Book struct {
	Title  string
	author string
	Pages  int
	Year   int
}

func (b Book) GetInfo() string {
	return fmt.Sprintf("Title: %s, Author: %s, Pages: %d, Year: %d", b.Title, b.author, b.Pages, b.Year)
}

func (b Book) IsLong() bool {
	return b.Pages > 300
}

func (b *Book) SetAuthor(name string) {
	b.author = name
}

func main() {

	fmt.Println(Book{Title: "Kill the bill"})
	fmt.Println(Book{Title: "The Great Gatsby", author: "F. Scott Fitzgerald", Pages: 180, Year: 1925})

	book := Book{Title: "1984", author: "George Orwell", Pages: 328, Year: 1948}
	fmt.Println(book.GetInfo())
	fmt.Println(book.IsLong())

	book.SetAuthor("Orwell")
	fmt.Println(book.GetInfo())

}
