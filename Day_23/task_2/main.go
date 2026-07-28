package main

import (
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type User struct {
	ID    int    `db:"id"`
	Name  string `db:"name"`
	Email string `db:"email"`
	Age   int    `db:"age"`
}

func main() {
	connStr := "user=ivan password=123 dbname=go_db host=localhost sslmode=disable"

	// sqlx.Connect вместо sql.Open + Ping
	db, err := sqlx.Connect("postgres", connStr)
	if err != nil {
		log.Fatal("Ошибка подключения:", err)
	}
	defer db.Close()

	fmt.Println("Подключение успешно!")

	// db.MustExec — падает при ошибке (не надо проверять err)
	db.MustExec(`
		CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			name VARCHAR(100) NOT NULL,
			email VARCHAR(100) UNIQUE NOT NULL,
			age INTEGER CHECK (age >= 0 AND age <= 150),
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)
	`)

	fmt.Println("Таблица users создана или уже существует!")

	// Добавление пользователя
	newUser := User{Name: "Ivan", Email: "ivan@mail.com", Age: 30}

	// db.MustExec — вставка
	db.MustExec(
		"INSERT INTO users (name, email, age) VALUES ($1, $2, $3)",
		newUser.Name, newUser.Email, newUser.Age,
	)
	fmt.Println("Пользователь добавлен!")

	// Получение всех пользователей
	var users []User
	err = db.Select(&users, "SELECT * FROM users")
	if err != nil {
		log.Fatal("Ошибка получения:", err)
	}
	fmt.Println("Все пользователи:", users)

	// Получение одного пользователя
	var user User
	err = db.Get(&user, "SELECT * FROM users WHERE name = $1", "Ivan")
	if err != nil {
		log.Fatal("Ошибка получения:", err)
	}
	fmt.Printf("Найден: %+v\n", user)

	// Обновление
	db.MustExec("UPDATE users SET age = $1 WHERE name = $2", 31, "Ivan")
	fmt.Println("Возраст обновлён!")

	// Проверка обновления
	db.Get(&user, "SELECT * FROM users WHERE name = $1", "Ivan")
	fmt.Printf("После обновления: %+v\n", user)

	// Удаление
	db.MustExec("DELETE FROM users WHERE name = $1", "Ivan")
	fmt.Println("Пользователь удалён!")

	// Проверка после удаления
	db.Select(&users, "SELECT * FROM users")
	fmt.Println("Осталось пользователей:", users)
}
