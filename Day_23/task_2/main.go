// доделать

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
	// 1️⃣ Подключение
	connStr := "user=ivan password=123 dbname=go_db host=localhost sslmode=disable"
	db, err := sqlx.Connect("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	fmt.Println("Подключено!")

	// 2️⃣ Создание таблицы
	db.MustExec(`
		CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			name TEXT,
			email TEXT UNIQUE,
			age INT
		)
	`)
	fmt.Println("Таблица готова!")

	// 3️⃣ INSERT
	newUser := User{Name: "Ivan", Email: "ivan@mail.com", Age: 30}
	db.MustExec("INSERT INTO users (name, email, age) VALUES ($1, $2, $3)",
		newUser.Name, newUser.Email, newUser.Age)
	fmt.Println("Пользователь добавлен!")

	// 4️⃣ SELECT все
	var users []User
	db.Select(&users, "SELECT * FROM users")
	fmt.Println("Все пользователи:", users)

	// 5️⃣ SELECT один
	var user User
	db.Get(&user, "SELECT * FROM users WHERE name=$1", "Ivan")
	fmt.Println("Найден:", user)

	// 6️⃣ UPDATE
	db.MustExec("UPDATE users SET age=$1 WHERE name=$2", 31, "Ivan")
	fmt.Println("Возраст обновлён!")

	// 7️⃣ DELETE
	db.MustExec("DELETE FROM users WHERE name=$1", "Ivan")
	fmt.Println("Пользователь удалён!")

	// 8️⃣ Проверка
	var remaining []User
	db.Select(&remaining, "SELECT * FROM users")
	fmt.Println("Осталось:", remaining)
}
