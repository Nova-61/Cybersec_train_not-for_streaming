package models

import (
	"time"
)

// User — то, что хранится в базе данных
type User struct {
	ID        int       `json:"id" db:"id"`
	Username  string    `json:"username" db:"username"`
	Email     string    `json:"email" db:"email"`
	Password  string    `json:"-" db:"password"` // "-" значит: никогда не выводить в JSON
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// RegisterRequest — то, что клиент присылает на /register
type RegisterRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// LoginRequest — то, что клиент присылает на /login
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// AuthResponse — то, что мы отдаём клиенту после логина/регистрации/refresh
type AuthResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
}

// RefreshRequest — то, что клиент присылает на /refresh
type RefreshRequest struct {
	RefreshToken string `json:"refresh_token"`
}

// Session — то, что мы храним в Redis. Не путать с моделью User —
// это отдельная сущность про конкретный "заход" пользователя в систему.
type Session struct {
	UserID       int       `json:"user_id"`
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	CreatedAt    time.Time `json:"created_at"`
	ExpiresAt    time.Time `json:"expires_at"`
}
