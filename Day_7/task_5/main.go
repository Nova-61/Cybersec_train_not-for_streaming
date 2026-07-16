package main

import (
	"encoding/json"
	"fmt"
	"myproject/models"
	"myproject/storage"
)

func main() {
	user := models.User{
		Name:  "Иван",
		Age:   30,
		Email: "ivan@example.com",
	}

	userData, err := json.Marshal(user)
	if err != nil {
		fmt.Println("Ошибка сериализации:", err)
		return
	}

	fileStorage, err := storage.NewFileStorage("data")
	if err != nil {
		fmt.Println("Ошибка создания FileStorage:", err)
		return
	}

	memoryStorage := storage.NewMemoryStorage()

	err = fileStorage.Save("user.json", userData)
	if err != nil {
		fmt.Println("Ошибка сохранения в файл:", err)
		return
	}
	fmt.Println("✅ Сохранено в FileStorage")

	err = memoryStorage.Save("user.json", userData)
	if err != nil {
		fmt.Println("Ошибка сохранения в память:", err)
		return
	}
	fmt.Println("✅ Сохранено в MemoryStorage")

	fmt.Println("\n--- Загрузка из FileStorage ---")
	fileData, err := fileStorage.Load("user.json")
	if err != nil {
		fmt.Println("Ошибка загрузки из файла:", err)
		return
	}
	var fileUser models.User
	err = json.Unmarshal(fileData, &fileUser)
	if err != nil {
		fmt.Println("Ошибка десериализации:", err)
		return
	}
	fmt.Println(fileUser)

	fmt.Println("\n--- Загрузка из MemoryStorage ---")
	memoryData, err := memoryStorage.Load("user.json")
	if err != nil {
		fmt.Println("Ошибка загрузки из памяти:", err)
		return
	}
	var memoryUser models.User
	err = json.Unmarshal(memoryData, &memoryUser)
	if err != nil {
		fmt.Println("Ошибка десериализации:", err)
		return
	}
	fmt.Println(memoryUser)
}
