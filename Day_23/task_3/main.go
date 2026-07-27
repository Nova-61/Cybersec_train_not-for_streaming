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

	db, err := sqlx.Connect("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	fmt.Println("Подключение успешно!")

	// Таблица users
	db.MustExec(`
		CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			name VARCHAR(100) NOT NULL,
			email VARCHAR(100) UNIQUE NOT NULL,
			age INTEGER
		)
	`)

	// Таблица logs
	db.MustExec(`
		CREATE TABLE IF NOT EXISTS logs (
			id SERIAL PRIMARY KEY,
			action TEXT,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)
	`)

	fmt.Println("Таблицы созданы!")

	// ========== CREATE USER WITH LOG ==========
	fmt.Println("\n=== CREATE USER WITH LOG ===")
	err = CreateUserWithLog(db, User{Name: "Ivan", Email: "ivan@mail.com", Age: 30}, "user created")
	if err != nil {
		fmt.Println("Ошибка:", err)
	} else {
		fmt.Println("Пользователь создан!")
	}

	// ========== UPDATE USER WITH LOG ==========
	fmt.Println("\n=== UPDATE USER WITH LOG ===")
	err = UpdateUserWithLog(db, 1, 31, "user updated")
	if err != nil {
		fmt.Println("Ошибка:", err)
	} else {
		fmt.Println("Пользователь обновлён!")
	}

	// ========== DELETE USER WITH LOG ==========
	fmt.Println("\n=== DELETE USER WITH LOG ===")
	err = DeleteUserWithLog(db, 1, "user deleted")
	if err != nil {
		fmt.Println("Ошибка:", err)
	} else {
		fmt.Println("Пользователь удалён!")
	}

	// ========== ПОКАЗАТЬ ВСЕХ ПОЛЬЗОВАТЕЛЕЙ ==========
	fmt.Println("\n=== ВСЕ ПОЛЬЗОВАТЕЛИ ===")
	var users []User
	db.Select(&users, "SELECT * FROM users")
	fmt.Println(users)

	// ========== ПОКАЗАТЬ ЛОГИ ==========
	fmt.Println("\n=== ЛОГИ ===")
	var logs []string
	db.Select(&logs, "SELECT action FROM logs")
	for _, l := range logs {
		fmt.Println("-", l)
	}
}

// ========== ТРАНЗАКЦИИ ==========

func CreateUserWithLog(db *sqlx.DB, user User, action string) error {
	// Начинаем транзакцию
	tx, err := db.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback() // откат, если не закоммитим

	// 1️⃣ Добавляем пользователя
	_, err = tx.Exec(
		"INSERT INTO users (name, email, age) VALUES ($1, $2, $3)",
		user.Name, user.Email, user.Age,
	)
	if err != nil {
		return err
	}

	// 2️⃣ Добавляем лог
	_, err = tx.Exec("INSERT INTO logs (action) VALUES ($1)", action)
	if err != nil {
		return err
	}

	// 3️⃣ Коммитим
	return tx.Commit()
}

func UpdateUserWithLog(db *sqlx.DB, id int, newAge int, action string) error {
	tx, err := db.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// 1️⃣ Обновляем пользователя
	_, err = tx.Exec("UPDATE users SET age = $1 WHERE id = $2", newAge, id)
	if err != nil {
		return err
	}

	// 2️⃣ Добавляем лог
	_, err = tx.Exec("INSERT INTO logs (action) VALUES ($1)", action)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func DeleteUserWithLog(db *sqlx.DB, id int, action string) error {
	tx, err := db.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// 1️⃣ Удаляем пользователя
	_, err = tx.Exec("DELETE FROM users WHERE id = $1", id)
	if err != nil {
		return err
	}

	// 2️⃣ Добавляем лог
	_, err = tx.Exec("INSERT INTO logs (action) VALUES ($1)", action)
	if err != nil {
		return err
	}

	return tx.Commit()
}
