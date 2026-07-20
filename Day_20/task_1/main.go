package main

import (
	"fmt"
	"os"
	"strings"
)

type FileManager struct {
	Filename string
}

func (fm FileManager) Write(content string) error {
	err := os.WriteFile(fm.Filename, []byte(content), 0644)
	if err != nil {
		return fmt.Errorf("ошибка записи: %v", err)
	}
	return nil
}

func (fm FileManager) Read() (string, error) {
	data, err := os.ReadFile(fm.Filename)
	if err != nil {
		return "", fmt.Errorf("ошибка чтения: %v", err)
	}
	return string(data), nil
}

func (fm FileManager) Action(action string, content ...string) error {
	switch action {
	case "write":
		if len(content) == 0 {
			return fmt.Errorf("для записи нужен контент")
		}
		fullContent := strings.Join(content, " ")
		return fm.Write(fullContent)

	case "read":
		data, err := fm.Read()
		if err != nil {
			return err
		}
		fmt.Println("Содержимое:", data)
		return nil

	case "write_read":
		if len(content) == 0 {
			return fmt.Errorf("для записи нужен контент")
		}
		fullContent := strings.Join(content, " ")
		err := fm.Write(fullContent)
		if err != nil {
			return err
		}
		data, err := fm.Read()
		if err != nil {
			return err
		}
		fmt.Println("Содержимое:", data)
		return nil

	default:
		return fmt.Errorf("неизвестное действие: %s", action)
	}
}

func main() {
	file := FileManager{Filename: "file.txt"}

	// Записываем одну строку
	err := file.Action("write", "Hello, World!")
	if err != nil {
		fmt.Println(err)
	}

	// Записываем несколько строк
	err = file.Action("write", "Hello", "World!", "Go")
	if err != nil {
		fmt.Println(err)
	}

	// Читаем
	err = file.Action("read")
	if err != nil {
		fmt.Println(err)
	}
}
