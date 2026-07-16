package main

import (
	"fmt"
	"myproject/models"
	"myproject/utils"
)

func init() {
	fmt.Println("Main package initialized")
}

func main() {
	user := models.User{
		Name:  "Иван",
		Email: "ivan@example.com",
	}

	_, err := user.SetAge(30)
	if err != nil {
		fmt.Println("Ошибка установки возраста:", err)
		return
	}

	result, err := utils.PrintUser(user)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println(result)
}
