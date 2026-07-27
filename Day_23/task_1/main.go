package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/lib/pq"
)

func connectDB() (*sql.DB, error) {
	connStr := "user=postgres password=secret dbname=go_db sslmode=disable"

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	fmt.Println("Подключено к БД!")
	return db, nil
}

func main() {
	db, err := connectDB()
	if err != nil {
		log.Fatal("Ошибка подключения:", err)
	}
	defer db.Close()

	fmt.Println("вроде все хорошо")
}
