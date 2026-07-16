package utils

import (
	"fmt"
	"myproject/models"
)

var AppVerrsion string

func PrintUser(u models.User) (string, error) {
	return fmt.Sprintf("Его email: %s, имя: %s", u.Email, u.Name), nil
}

func init() {
	fmt.Println("Utils package initialized")
	AppVerrsion = "1.0.0"
}
