package middleware

import (
	"context"
	"net/http"
	"strings"

	"auth/internal/utils"
)

// contextKey — приватный тип, а не string. Это гарантирует, что НИКАКОЙ другой
// пакет не сможет случайно (или намеренно) подсунуть значение по этому же ключу,
// даже если тоже назовёт переменную "user_id" — сравнение типов в Go учитывает
// сам тип ключа, не только его значение.
type contextKey int

const (
	// UserIDContextKey и UsernameContextKey экспортированы (с большой буквы),
	// потому что handler-пакет тоже должен уметь их прочитать -
	// middleware кладёт значения, handler их читает.
	UserIDContextKey contextKey = iota
	UsernameContextKey
)

type AuthMiddleware struct {
	jwtManager *utils.JWTManager
}

func NewAuthMiddleware(jwtManager *utils.JWTManager) *AuthMiddleware {
	return &AuthMiddleware{jwtManager: jwtManager}
}

// Validate оборачивает защищённый хендлер: проверяет Bearer-токен и,
// если всё ок, прокидывает user_id/username дальше через контекст запроса.
func (m *AuthMiddleware) Validate(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header required", http.StatusUnauthorized)
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, "Invalid authorization header format. Use: Bearer <token>", http.StatusUnauthorized)
			return
		}
		token := parts[1]

		claims, err := m.jwtManager.ValidateToken(token)
		if err != nil {
			// Намеренно НЕ включаем err.Error() в ответ клиенту (см. обсуждение
			// про логи vs HTTP-ответы) — общий "Invalid token" достаточно
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), UserIDContextKey, claims.UserID)
		ctx = context.WithValue(ctx, UsernameContextKey, claims.Username)

		next(w, r.WithContext(ctx))
	}
}