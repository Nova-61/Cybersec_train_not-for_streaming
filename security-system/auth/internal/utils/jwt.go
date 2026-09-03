package utils

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Claims — что мы кладём внутрь access-токена
type Claims struct {
	UserID   int    `json:"user_id"`
	Username string `json:"username"`
	Type     string `json:"type"` // всегда "access" — защита от подмены типа токена
	jwt.RegisteredClaims
}

// RefreshClaims — что кладём внутрь refresh-токена.
// Намеренно НЕ содержит username — refresh-токену не нужны лишние данные,
// его единственная задача — доказать право получить новый access-токен.
type RefreshClaims struct {
	UserID int    `json:"user_id"`
	Type   string `json:"type"` // всегда "refresh"
	jwt.RegisteredClaims
}

type JWTManager struct {
	secretKey     string
	tokenExpiry   time.Duration
	refreshExpiry time.Duration
}

func NewJWTManager(secretKey string, tokenExpiry, refreshExpiry time.Duration) *JWTManager {
	return &JWTManager{
		secretKey:     secretKey,
		tokenExpiry:   tokenExpiry,
		refreshExpiry: refreshExpiry,
	}
}

// GenerateAccessToken создаёт короткоживущий токен с данными пользователя
func (j *JWTManager) GenerateAccessToken(userID int, username string) (string, error) {
	claims := Claims{
		UserID:   userID,
		Username: username,
		Type:     "access",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(j.tokenExpiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(j.secretKey))
}

// GenerateRefreshToken создаёт долгоживущий токен с user_id, чтобы потом
// можно было выдать новый access-токен, не заставляя пользователя логиниться заново
func (j *JWTManager) GenerateRefreshToken(userID int) (string, error) {
	claims := RefreshClaims{
		UserID: userID,
		Type:   "refresh",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(j.refreshExpiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(j.secretKey))
}

// ValidateToken проверяет подпись и срок действия access-токена
func (j *JWTManager) ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(j.secretKey), nil
	})
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}
	// Проверяем, что это именно access-токен, а не подсунутый refresh-токен
	if claims.Type != "access" {
		return nil, errors.New("wrong token type")
	}
	return claims, nil
}

// ValidateRefreshToken проверяет refresh-токен и достаёт из него user_id
func (j *JWTManager) ValidateRefreshToken(tokenString string) (int, error) {
	token, err := jwt.ParseWithClaims(tokenString, &RefreshClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(j.secretKey), nil
	})
	if err != nil {
		return 0, errors.New("invalid refresh token")
	}

	claims, ok := token.Claims.(*RefreshClaims)
	if !ok || !token.Valid {
		return 0, errors.New("invalid refresh token")
	}
	// Проверяем, что это именно refresh-токен, а не подсунутый access-токен
	if claims.Type != "refresh" {
		return 0, errors.New("wrong token type")
	}

	return claims.UserID, nil
}
