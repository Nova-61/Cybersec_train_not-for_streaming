package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"auth/internal/config"
	"auth/internal/middleware"
	"auth/internal/models"
	"auth/internal/repository"
	"auth/internal/utils"
)

type AuthHandler struct {
	userRepo    *repository.UserRepository
	sessionRepo *repository.SessionRepository
	jwtManager  *utils.JWTManager
	config      *config.Config
}

func NewAuthHandler(
	userRepo *repository.UserRepository,
	sessionRepo *repository.SessionRepository,
	jwtManager *utils.JWTManager,
	cfg *config.Config,
) *AuthHandler {
	return &AuthHandler{
		userRepo:    userRepo,
		sessionRepo: sessionRepo,
		jwtManager:  jwtManager,
		config:      cfg,
	}
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req models.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Базовая валидация
	if req.Username == "" || req.Email == "" || req.Password == "" {
		http.Error(w, "All fields are required", http.StatusBadRequest)
		return
	}
	if len(req.Password) < 6 {
		http.Error(w, "Password must be at least 6 characters", http.StatusBadRequest)
		return
	}

	// Проверка уникальности. err игнорируем намеренно здесь -
	// если БД реально недоступна, следующий вызов Create всё равно провалится с 500
	existing, _ := h.userRepo.GetByUsername(req.Username)
	if existing != nil {
		http.Error(w, "Username already exists", http.StatusConflict)
		return
	}
	existing, _ = h.userRepo.GetByEmail(req.Email)
	if existing != nil {
		http.Error(w, "Email already exists", http.StatusConflict)
		return
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	user := &models.User{
		Username: req.Username,
		Email:    req.Email,
		Password: hashedPassword,
	}
	if err := h.userRepo.Create(user); err != nil {
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	h.issueTokensAndRespond(w, r.Context(), user, http.StatusCreated)
}

// issueTokensAndRespond — общая логика для Register и Login:
// генерирует пару токенов, сохраняет сессию в Redis, отвечает клиенту.
func (h *AuthHandler) issueTokensAndRespond(w http.ResponseWriter, ctx context.Context, user *models.User, statusCode int) {
	accessToken, err := h.jwtManager.GenerateAccessToken(user.ID, user.Username)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}
	refreshToken, err := h.jwtManager.GenerateRefreshToken(user.ID)
	if err != nil {
		http.Error(w, "Failed to generate refresh token", http.StatusInternalServerError)
		return
	}

	session := &models.Session{
		UserID:       user.ID,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		CreatedAt:    time.Now(),
		ExpiresAt:    time.Now().Add(time.Hour * time.Duration(h.config.JWTExpiry)),
	}
	if err := h.sessionRepo.Save(ctx, user.ID, session); err != nil {
		http.Error(w, "Failed to save session", http.StatusInternalServerError)
		return
	}

	response := models.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    h.config.JWTExpiry * 3600,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(response)
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req models.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	if req.Username == "" || req.Password == "" {
		http.Error(w, "Username and password are required", http.StatusBadRequest)
		return
	}

	user, err := h.userRepo.GetByUsername(req.Username)
	if err != nil {
		// Реальная ошибка инфраструктуры (БД недоступна) — это НЕ "неверный пароль"
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	if user == nil {
		// Пользователя действительно нет — вот это уже 401
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	if !utils.CheckPasswordHash(req.Password, user.Password) {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	h.issueTokensAndRespond(w, r.Context(), user, http.StatusOK)
}

func (h *AuthHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	var req models.RefreshRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	if req.RefreshToken == "" {
		http.Error(w, "Refresh token is required", http.StatusBadRequest)
		return
	}

	// 1. Проверяем подпись, срок действия и тип токена (см. utils/jwt.go)
	userID, err := h.jwtManager.ValidateRefreshToken(req.RefreshToken)
	if err != nil {
		http.Error(w, "Invalid refresh token", http.StatusUnauthorized)
		return
	}

	// 2. Пользователь мог быть удалён после выдачи токена — перепроверяем
	user, err := h.userRepo.GetByID(userID)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	if user == nil {
		http.Error(w, "User not found", http.StatusUnauthorized)
		return
	}

	// 3. Ключевая проверка: токен должен совпадать с тем, что лежит в Redis.
	// Если сессия была удалена (logout) или уже обновлена другим запросом —
	// этот refresh-токен больше не действителен, даже если подпись валидна.
	ctx := r.Context()
	session, err := h.sessionRepo.Get(ctx, userID)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	if session == nil || session.RefreshToken != req.RefreshToken {
		http.Error(w, "Invalid refresh token", http.StatusUnauthorized)
		return
	}

	// 4. Всё ок — выдаём новую пару токенов (refresh token rotation)
	h.issueTokensAndRespond(w, ctx, user, http.StatusOK)
}

// Logout удаляет сессию из Redis. См. обсуждение выше про ограничения этого подхода:
// уже выданный access-токен остаётся валидным по подписи до истечения exp,
// но обновить его через /refresh после этого будет уже нельзя.
func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDContextKey).(int)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	if err := h.sessionRepo.Delete(r.Context(), userID); err != nil {
		http.Error(w, "Failed to logout", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// Validate — эндпоинт для других сервисов/клиента, чтобы проверить,
// что access-токен ещё действителен и получить актуальные данные о пользователе
func (h *AuthHandler) Validate(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDContextKey).(int)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	user, err := h.userRepo.GetByID(userID)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	if user == nil {
		http.Error(w, "User not found", http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"valid":    true,
		"user_id":  user.ID,
		"username": user.Username,
	})
}