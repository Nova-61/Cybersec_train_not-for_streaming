package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"auth/internal/config"
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

// issueTokensAndRespond - общая логика для Register и Login:
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

//Честно говоря, я не уверен, что это правильный способ сделать logout. Я просто удаляю сессию из Redis, и всё. Но если кто-то ещё где-то хранит refresh-токен, он всё равно сможет получить новый access-токен. Поэтому в реальной жизни нужно ещё как-то инвалидировать refresh-токен, например, хранить их в БД и помечать как "использованный" или "отозванный". Но для учебного проекта этого достаточно.
