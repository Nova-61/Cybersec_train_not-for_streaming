package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

func main() {
	connStr := "user=ivan password=123 dbname=go_db host=localhost sslmode=disable"

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Ошибка подключения:", err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatal("Ошибка проверки соединения:", err)
	}

	fmt.Println("Подключение успешно!")

	// Создание таблицы
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			name VARCHAR(100) NOT NULL,
			email VARCHAR(100) UNIQUE NOT NULL,
			age INTEGER CHECK (age >= 0 AND age <= 150),
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		log.Fatal("Ошибка создания таблицы:", err)
	}
	fmt.Println("Таблица users создана или уже существует!")

	// Добавление 3 пользователей
	fmt.Println("\n--- Добавление пользователей ---")

	users := []struct {
		name  string
		email string
		age   int
	}{
		{"Ivan", "ivan@mail.com", 30},
		{"Petr", "petr@mail.com", 25},
		{"Anna", "anna@mail.com", 28},
	}

	for _, u := range users {
		_, err := db.Exec(
			"INSERT INTO users (name, email, age) VALUES ($1, $2, $3)",
			u.name, u.email, u.age,
		)
		if err != nil {
			log.Printf("Ошибка добавления пользователя %s: %v\n", u.name, err)
		} else {
			fmt.Printf("Добавлен пользователь: %s (%d лет)\n", u.name, u.age)
		}
	}

	// Выбор всех пользователей
	fmt.Println("\n--- Все пользователи ---")
	rows, err := db.Query("SELECT id, name, email, age, created_at FROM users")
	if err != nil {
		log.Fatal("Ошибка запроса:", err)
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		var name, email string
		var age int
		var createdAt string
		err := rows.Scan(&id, &name, &email, &age, &createdAt)
		if err != nil {
			log.Fatal("Ошибка чтения:", err)
		}
		fmt.Printf("ID: %d, Имя: %s, Email: %s, Возраст: %d, Создан: %s\n",
			id, name, email, age, createdAt[:19])
	}

	// Обновление возраста
	fmt.Println("\n--- Обновление возраста ---")
	result, err := db.Exec("UPDATE users SET age = $1 WHERE name = $2", 31, "Ivan")
	if err != nil {
		log.Fatal("Ошибка обновления:", err)
	}
	rowsAffected, _ := result.RowsAffected()
	fmt.Printf("Обновлено строк: %d\n", rowsAffected)

	// Проверяем обновление
	rows, err = db.Query("SELECT name, age FROM users WHERE name = 'Ivan'")
	if err != nil {
		log.Fatal("Ошибка запроса:", err)
	}
	defer rows.Close()

	if rows.Next() {
		var name string
		var age int
		rows.Scan(&name, &age)
		fmt.Printf("После обновления: %s, возраст: %d\n", name, age)
	}

	// Удаление пользователя
	fmt.Println("\n--- Удаление пользователя ---")
	result, err = db.Exec("DELETE FROM users WHERE name = $1", "Petr")
	if err != nil {
		log.Fatal("Ошибка удаления:", err)
	}
	rowsAffected, _ = result.RowsAffected()
	fmt.Printf("Удалено строк: %d\n", rowsAffected)

	// Проверяем всех оставшихся пользователей
	fmt.Println("\n--- Оставшиеся пользователи ---")
	rows, err = db.Query("SELECT id, name, email, age FROM users")
	if err != nil {
		log.Fatal("Ошибка запроса:", err)
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		var name, email string
		var age int
		rows.Scan(&id, &name, &email, &age)
		fmt.Printf("ID: %d, Имя: %s, Email: %s, Возраст: %d\n", id, name, email, age)
	}

	fmt.Println("\nПрограмма завершена!")
}
