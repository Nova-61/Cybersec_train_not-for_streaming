package models

import (
	"time"
)

type User struct { // Модель пользователя, которая представляет собой структуру данных, содержащую информацию о пользователе в системе.
	ID        int       `json:"id" db:"id"`
	Username  string    `json:"username" db:"username"`
	Password  string    `json:"-" db:"password"` // Никогда не возвращай пароль в ответе на запрос, поэтому используем тег json:"-" чтобы скрыть его при сериализации в JSON.
	Email     string    `json:"email" db:"email"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

type RegisterRequest struct { // Запрос на регистрацию нового пользователя
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

type LoginRequest struct { // Запрос на аутентификацию пользователя
	Username string `json:"username"`
	Password string `json:"password"`
}

type AuthResponse struct { // Отдаем пользователю токены после успешной аутентификации (Это типо как билет в кинотеатр, который подтверждает, что пользователь прошел проверку и может получить доступ к защищенным ресурсам)
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int64  `json:"expires_in"`
}
