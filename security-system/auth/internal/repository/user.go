package repository

import (
	"database/sql"
	"fmt"
	"log"

	"auth/internal/models"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

// Create добавляет пользователя и заполняет ID/CreatedAt/UpdatedAt из ответа БД
func (r *UserRepository) Create(user *models.User) error {
	query := `INSERT INTO users (username, email, password) VALUES ($1, $2, $3) RETURNING id, created_at, updated_at`
	err := r.db.QueryRow(query, user.Username, user.Email, user.Password).Scan(
		&user.ID, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		log.Printf("Error creating user: %v", err)
		return err
	}
	return nil
}

// Все последующие функции почти одинаковые, их задача просто по единственоуму параметру доставать и сверять пользователя с базой
// по никнейму, почте или айди. Если пользователь не найден, возвращаем nil, nil. Если произошла ошибка — nil, err. Если найден — &user, nil.
func (r *UserRepository) GetByUsername(username string) (*models.User, error) {
	var user models.User
	query := `SELECT id, username, email, password, created_at, updated_at FROM users WHERE username = $1`
	err := r.db.QueryRow(query, username).Scan(
		&user.ID, &user.Username, &user.Email, &user.Password, &user.CreatedAt, &user.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		log.Printf("Error getting user: %v", err)
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) GetByEmail(email string) (*models.User, error) {
	var user models.User
	query := `SELECT id, username, email, password, created_at, updated_at FROM users WHERE email = $1`
	err := r.db.QueryRow(query, email).Scan(
		&user.ID, &user.Username, &user.Email, &user.Password, &user.CreatedAt, &user.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		log.Printf("Error getting user: %v", err)
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) GetByID(id int) (*models.User, error) {
	var user models.User
	query := `SELECT id, username, email, password, created_at, updated_at FROM users WHERE id = $1`
	err := r.db.QueryRow(query, id).Scan(
		&user.ID, &user.Username, &user.Email, &user.Password, &user.CreatedAt, &user.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		log.Printf("Error getting user: %v", err)
		return nil, err
	}
	return &user, nil
}

// CreateTable создает таблицу users, если она НЕ СУЩЕСТВЕТ, и добавляет индексы на username и email
func (r *UserRepository) CreateTable() error {
	query := `
	CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		username VARCHAR(50) UNIQUE NOT NULL,
		email VARCHAR(100) UNIQUE NOT NULL,
		password VARCHAR(255) NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);
	CREATE INDEX IF NOT EXISTS idx_users_username ON users(username);
	CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
	`
	_, err := r.db.Exec(query)
	if err != nil {
		log.Printf("Error creating users table: %v", err)
		return err
	}
	fmt.Println("Users table created")
	return nil
}
