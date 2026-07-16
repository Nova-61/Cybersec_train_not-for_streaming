package main

import (
	"fmt"
	"myproject/models"
	"myproject/utils"
)

func main() {
	user := models.User{
		Name:  "Иван",
		Age:   "30",
		Email: "ivan@example.com",
	}
	result, err := utils.PrintUser(user)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println(result)
}
