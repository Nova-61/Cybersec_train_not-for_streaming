package models

import "fmt"

type User struct {
	Name  string
	Age   int
	Email string
}

func (u User) String() string {
	return fmt.Sprintf("Name: %s, Age: %d, Email: %s", u.Name, u.Age, u.Email)
}
