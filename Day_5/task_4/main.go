package main

import (
	"fmt"
	"strings"
)

type Author struct {
	Name      string
	BirthYear int
}

type Book struct {
	Title string
	Author
	Pages     int
	Year      int
	Publisher string
}

type User struct {
	Username string
	Email    string
	Age      int
	Password string
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

func (u User) Validate() error {
	if u.Username == "" {
		return fmt.Errorf("username is required")
	}

	if !strings.Contains(u.Email, "@") {
		return fmt.Errorf("email must contain '@'")
	}

	if u.Age < 18 || u.Age > 100 {
		return fmt.Errorf("age must be between 18 and 100")
	}

	if len(u.Password) < 8 {
		return fmt.Errorf("password must be at least 8 characters")
	}

	return nil
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

	fmt.Println("\n--- User Validation ---")

	user := User{
		Username: "john_doe",
		Email:    "john@example.com",
		Age:      25,
		Password: "securepass",
	}

	if err := user.Validate(); err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("User is valid")
	}

	invalidUser := User{
		Username: "",
		Email:    "invalid",
		Age:      15,
		Password: "123",
	}

	if err := invalidUser.Validate(); err != nil {
		fmt.Println("Error:", err)
	}
}
