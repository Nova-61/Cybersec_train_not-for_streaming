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

type Log struct {
	ID     int    `db:"id"`
	Action string `db:"action"`
}

func main() {
	connStr := "user=ivan password=123 dbname=go_db host=localhost sslmode=disable"
	db, err := sqlx.Connect("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	fmt.Println("Подключение успешно!")

	// Создаём таблицы
	db.MustExec(`
		CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			name VARCHAR(100) NOT NULL,
			email VARCHAR(100) UNIQUE NOT NULL,
			age INTEGER
		)
	`)

	db.MustExec(`
		CREATE TABLE IF NOT EXISTS logs (
			id SERIAL PRIMARY KEY,
			action TEXT,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)
	`)

	fmt.Println("Таблицы созданы!")

	// 1. CREATE USER (NamedExec)
	fmt.Println("\n=== CREATE USER ===")
	user := User{Name: "Ivan", Email: "ivan@mail.com", Age: 30}
	err = CreateUser(db, user)
	if err != nil {
		fmt.Println("Ошибка:", err)
	} else {
		fmt.Println("Пользователь создан!")
	}

	// 2. UPDATE USER (NamedExec)
	fmt.Println("\n=== UPDATE USER ===")
	userUpdate := User{ID: 1, Name: "Ivan Petrov", Email: "ivan_new@mail.com", Age: 31}
	err = UpdateUser(db, userUpdate)
	if err != nil {
		fmt.Println("Ошибка:", err)
	} else {
		fmt.Println("Пользователь обновлён!")
	}

	// 3. CREATE USER WITH LOG (NamedExec)
	fmt.Println("\n=== CREATE USER WITH LOG ===")
	err = CreateUserWithLog(db, User{Name: "Anna", Email: "anna@mail.com", Age: 28}, "user created")
	if err != nil {
		fmt.Println("Ошибка:", err)
	} else {
		fmt.Println("Пользователь с логом создан!")
	}

	//  ПОКАЗАТЬ ВСЕХ ПОЛЬЗОВАТЕЛЕЙ
	fmt.Println("\n=== ВСЕ ПОЛЬЗОВАТЕЛИ ===")
	var users []User
	db.Select(&users, "SELECT * FROM users")
	for _, u := range users {
		fmt.Printf("ID: %d, Имя: %s, Email: %s, Возраст: %d\n",
			u.ID, u.Name, u.Email, u.Age)
	}

	//  ПОКАЗАТЬ ЛОГИ
	fmt.Println("\n=== ЛОГИ ===")
	var logs []Log
	db.Select(&logs, "SELECT * FROM logs")
	for _, l := range logs {
		fmt.Printf("ID: %d, Действие: %s\n", l.ID, l.Action)
	}
}

// CREATE USER через NamedExec
func CreateUser(db *sqlx.DB, user User) error {
	query := `INSERT INTO users (name, email, age) VALUES (:name, :email, :age)`
	_, err := db.NamedExec(query, user)
	return err
}

// UPDATE USER через NamedExec
func UpdateUser(db *sqlx.DB, user User) error {
	query := `UPDATE users SET name = :name, email = :email, age = :age WHERE id = :id`
	_, err := db.NamedExec(query, user)
	return err
}

// CREATE USER WITH LOG через NamedExec (ТРАНЗАКЦИЯ)
func CreateUserWithLog(db *sqlx.DB, user User, action string) error {
	// Начинаем транзакцию
	tx, err := db.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Добавляем пользователя (NamedExec)
	userQuery := `INSERT INTO users (name, email, age) VALUES (:name, :email, :age)`
	_, err = tx.NamedExec(userQuery, user)
	if err != nil {
		return err
	}

	// Добавляем лог (тоже NamedExec, но для простоты используем map)
	logData := map[string]interface{}{
		"action": action,
	}
	logQuery := `INSERT INTO logs (action) VALUES (:action)`
	_, err = tx.NamedExec(logQuery, logData)
	if err != nil {
		return err
	}

	// Коммитим
	return tx.Commit()
}
