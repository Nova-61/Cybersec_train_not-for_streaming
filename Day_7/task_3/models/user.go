package models

import "fmt"

type User struct {
	Name  string
	age   int
	Email string
}

func init() {
	fmt.Println("Models package initialized")
}

func (u User) GetAge() (int, error) {
	if u.age == 0 {
		return 0, fmt.Errorf("age is not set")
	}
	if u.age < 0 || u.age > 150 {
		return 0, fmt.Errorf("invalid age: %d", u.age)
	}
	return u.age, nil
}

func (u *User) SetAge(age int) (int, error) {
	if age < 0 || age > 150 {
		return 0, fmt.Errorf("age cannot be negative or greater than 150")
	}
	u.age = age
	return u.age, nil
}
