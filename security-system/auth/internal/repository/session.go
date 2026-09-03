package repository

import (
	"context" // используется для передачи контекста выполнения, например, для таймаутов и отмены операций.
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"

	"auth/internal/models"
)

type SessionRepository struct {
	client *redis.Client
}

func NewSessionRepository(client *redis.Client) *SessionRepository {
	return &SessionRepository{client: client}
}

// Save сохраняет сессию с автоматическим истечением по TTL
func (r *SessionRepository) Save(ctx context.Context, userID int, session *models.Session) error {
	key := fmt.Sprintf("session:%d", userID) // Формируем ключ для Redis, чтобы хранить сессию конкретного пользователя. Например, "session:123" для пользователя с ID 123.
	data, err := json.Marshal(session)       // Преобразуем структуру сессии в JSON для хранения в Redis
	// session это структура, а redis хранит только строки, поэтому мы сериализуем её в JSON
	if err != nil {
		return err
	}

	expiry := time.Until(session.ExpiresAt)           // TTL для Redis, чтобы сессия автоматически удалялась
	return r.client.Set(ctx, key, data, expiry).Err() // Сохраняем сессию в Redis с TTL
}

// Get достаёт сессию. Возвращает nil, nil если сессии нет (не ошибка!)
func (r *SessionRepository) Get(ctx context.Context, userID int) (*models.Session, error) {
	key := fmt.Sprintf("session:%d", userID)
	data, err := r.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	var session models.Session                                     // Создаем переменную для хранения десериализованной сессии
	if err := json.Unmarshal([]byte(data), &session); err != nil { // Десериализуем JSON обратно в структуру Session. Если произошла ошибка, возвращаем её.
		return nil, err
	}

	return &session, nil
}

// Delete удаляет сессию — используется при logout
func (r *SessionRepository) Delete(ctx context.Context, userID int) error {
	key := fmt.Sprintf("session:%d", userID)
	return r.client.Del(ctx, key).Err()
}

// честно и добросовестно спиздил этот файл, у ии, becouse i dont know how to write redis code, i hope it works, and i will test it later
// and for my exployer, i swear that i will learn how to write redis code, and i will not use this code in production
// So, I think I understand what's going on here, but I wouldn't say I'm 100 percent sure. Anyway, don't worry about it; you'll have to deal with this later anyway, so I'll study it a little later than I do now. :)
// Я пиш на английском потому что устал постояно свапаться на русский
