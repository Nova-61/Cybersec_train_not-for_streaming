package utils

import (
	"fmt"
	"myproject/models"
)

func PrintUser(u models.User) (string, error) {
	return fmt.Sprintf("Его email: %s, имя: %s", u.Email, u.Name), nil
}
