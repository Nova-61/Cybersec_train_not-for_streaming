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

// ✅ 1. GET USERS BY AGE (WHERE age >= $1 AND age <= $2)
func GetUsersByAge(db *sqlx.DB, minAge, maxAge int) ([]User, error) {
	var users []User
	query := "SELECT * FROM users WHERE age >= $1 AND age <= $2"
	err := db.Select(&users, query, minAge, maxAge)
	return users, err
}

// ✅ 2. GET USER BY EMAIL (WHERE email = $1)
func GetUserByEmail(db *sqlx.DB, email string) (*User, error) {
	var user User
	query := "SELECT * FROM users WHERE email = $1"
	err := db.Get(&user, query, email)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// ✅ 3. GET USERS BY NAME (WHERE name = $1)
func GetUsersByName(db *sqlx.DB, name string) ([]User, error) {
	var users []User
	query := "SELECT * FROM users WHERE name = $1"
	err := db.Select(&users, query, name)
	return users, err
}

// ✅ 4. UPDATE USER EMAIL (UPDATE users SET email = $1 WHERE id = $2)
func UpdateUserEmail(db *sqlx.DB, id int, email string) error {
	query := "UPDATE users SET email = $1 WHERE id = $2"
	_, err := db.Exec(query, email, id)
	return err
}

func main() {
	connStr := "user=ivan password=123 dbname=go_db host=localhost sslmode=disable"

	db, err := sqlx.Connect("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	fmt.Println("Подключение успешно!")

	// Создаём таблицу
	db.MustExec(`
		CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			name VARCHAR(100) NOT NULL,
			email VARCHAR(100) UNIQUE NOT NULL,
			age INTEGER
		)
	`)

	// Добавляем тестовых пользователей
	db.MustExec("INSERT INTO users (name, email, age) VALUES ($1, $2, $3) ON CONFLICT DO NOTHING",
		"Ivan", "ivan@mail.com", 30)
	db.MustExec("INSERT INTO users (name, email, age) VALUES ($1, $2, $3) ON CONFLICT DO NOTHING",
		"Petr", "petr@mail.com", 25)
	db.MustExec("INSERT INTO users (name, email, age) VALUES ($1, $2, $3) ON CONFLICT DO NOTHING",
		"Anna", "anna@mail.com", 35)

	// ========== 1. GET USERS BY AGE ==========
	fmt.Println("\n=== ПОЛЬЗОВАТЕЛИ С ВОЗРАСТОМ 25-35 ===")
	users, err := GetUsersByAge(db, 25, 35)
	if err != nil {
		fmt.Println("Ошибка:", err)
	} else {
		for _, u := range users {
			fmt.Printf("ID: %d, Имя: %s, Email: %s, Возраст: %d\n",
				u.ID, u.Name, u.Email, u.Age)
		}
	}

	// ========== 2. GET USER BY EMAIL ==========
	fmt.Println("\n=== ПОЛЬЗОВАТЕЛЬ ПО EMAIL ===")
	user, err := GetUserByEmail(db, "ivan@mail.com")
	if err != nil {
		fmt.Println("Ошибка:", err)
	} else {
		fmt.Printf("Найден: ID=%d, Имя=%s, Email=%s, Возраст=%d\n",
			user.ID, user.Name, user.Email, user.Age)
	}

	// ========== 3. GET USERS BY NAME ==========
	fmt.Println("\n=== ПОЛЬЗОВАТЕЛИ ПО ИМЕНИ ===")
	usersByName, err := GetUsersByName(db, "Petr")
	if err != nil {
		fmt.Println("Ошибка:", err)
	} else {
		for _, u := range usersByName {
			fmt.Printf("ID: %d, Имя: %s, Email: %s, Возраст: %d\n",
				u.ID, u.Name, u.Email, u.Age)
		}
	}

	// ========== 4. UPDATE USER EMAIL ==========
	fmt.Println("\n=== ОБНОВЛЕНИЕ EMAIL ===")
	err = UpdateUserEmail(db, 1, "ivan_new@mail.com")
	if err != nil {
		fmt.Println("Ошибка:", err)
	} else {
		fmt.Println("Email обновлён!")
	}

	// Проверяем обновление
	updated, _ := GetUserByEmail(db, "ivan_new@mail.com")
	if updated != nil {
		fmt.Printf("Проверка: ID=%d, Имя=%s, Email=%s, Возраст=%d\n",
			updated.ID, updated.Name, updated.Email, updated.Age)
	}

	fmt.Println("\nГотово!")
}
