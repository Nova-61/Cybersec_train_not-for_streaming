package utils

import (
	"fmt"
	"myproject/models"
)

func PrintUser(u models.User) (string, error) {
	if u.Age == "" {
		return "", fmt.Errorf("age is empty")
	}
	return fmt.Sprintf("Ему %s лет, его email: %s, имя: %s", u.Age, u.Email, u.Name), nil
}
