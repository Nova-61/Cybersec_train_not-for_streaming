package main

import "fmt"

type Author struct {
	Name      string
	BirthYear int
}

type Book struct {
	Title     string
	Author    // встраивание Author
	Pages     int
	Year      int
	Publisher string
}

func (b Book) GetInfo() string {
	return fmt.Sprintf("Title: %s, Author: %s, Pages: %d, Year: %d, Publisher: %s",
		b.Title, b.Author.Name, b.Pages, b.Year, b.Publisher)
}

func (b Book) IsLong() bool {
	return b.Pages > 300
}

func (b *Book) SetAuthorName(name string) {
	b.Author.Name = name
}

func (b *Book) SetPublisher(publisher string) {
	b.Publisher = publisher
}

func main() {
	book := Book{
		Title: "1984",
		Author: Author{
			Name:      "George Orwell",
			BirthYear: 1903,
		},
		Pages:     328,
		Year:      1948,
		Publisher: "Secker & Warburg",
	}

	fmt.Println(book.GetInfo())
	fmt.Println(book.IsLong())

	book.SetAuthorName("Eric Arthur Blair")
	book.SetPublisher("Penguin Books")

	fmt.Println(book.GetInfo())
	fmt.Println("Author Name:", book.Author.Name)
	fmt.Println("Author Birth Year:", book.BirthYear)
}
