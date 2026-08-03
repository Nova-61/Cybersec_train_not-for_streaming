package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	// Определяем флаги
	filePath := flag.String("file", "logs.txt", "путь к файлу логов")
	fromDate := flag.String("from", "", "начальная дата (YYYY-MM-DD)")
	toDate := flag.String("to", "", "конечная дата (YYYY-MM-DD)")
	showStats := flag.Bool("stats", false, "показать статистику")

	// Парсим флаги
	flag.Parse()

	// Выводим значения флагов
	fmt.Println("=== Параметры запуска ===")
	fmt.Println("Файл:", *filePath)
	fmt.Println("От:", *fromDate)
	fmt.Println("До:", *toDate)
	fmt.Println("Статистика:", *showStats)

	// Проверяем, существует ли файл
	if _, err := os.Stat(*filePath); os.IsNotExist(err) {
		fmt.Printf("Ошибка: файл '%s' не существует!\n", *filePath)
		os.Exit(1)
	}

	fmt.Println("Файл существует!")
}
