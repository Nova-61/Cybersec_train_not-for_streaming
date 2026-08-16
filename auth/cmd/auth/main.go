package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-redis/redis/v8"
	_ "github.com/lib/pq"

	"auth/internal/config"
	"auth/internal/handler"
	"auth/internal/middleware"
	"auth/internal/repository"
	"auth/internal/utils"
)

func main() {
	// 1. КОНФИГ — падаем сразу, если JWT_SECRET не задан или слишком короткий
	cfg, err := config.New()
	if err != nil {
		log.Fatal("Invalid configuration: ", err)
	}

	// 2. ПОДКЛЮЧЕНИЕ К POSTGRES
	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName,
	)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Failed to connect to database: ", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatal("Failed to ping database: ", err)
	}
	fmt.Println("Connected to PostgreSQL")

	// 3. ПОДКЛЮЧЕНИЕ К REDIS
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisAddr,
		Password: cfg.RedisPassword,
		DB:       cfg.RedisDB,
	})
	if err := rdb.Ping(context.Background()).Err(); err != nil {
		log.Fatal("Failed to connect to Redis: ", err)
	}
	fmt.Println("Connected to Redis")
	defer rdb.Close()

	// 4. СБОРКА ЗАВИСИМОСТЕЙ (dependency injection руками)
	userRepo := repository.NewUserRepository(db)
	sessionRepo := repository.NewSessionRepository(rdb)
	jwtManager := utils.NewJWTManager(
		cfg.JWTSecret,
		time.Hour*time.Duration(cfg.JWTExpiry),
		time.Hour*24*time.Duration(cfg.RefreshExpiry),
	)
	authHandler := handler.NewAuthHandler(userRepo, sessionRepo, jwtManager, cfg)
	authMiddleware := middleware.NewAuthMiddleware(jwtManager)

	// 5. "МИГРАЦИЯ" — создание таблицы, если её ещё нет (см. обсуждение шага 5)
	if err := userRepo.CreateTable(); err != nil {
		log.Fatal("Failed to create tables: ", err)
	}

	// 6. РЕГИСТРАЦИЯ РОУТОВ
	mux := http.NewServeMux()
	mux.HandleFunc("/register", authHandler.Register)
	mux.HandleFunc("/login", authHandler.Login)
	mux.HandleFunc("/refresh", authHandler.Refresh)
	mux.HandleFunc("/logout", authMiddleware.Validate(authHandler.Logout))
	mux.HandleFunc("/validate", authMiddleware.Validate(authHandler.Validate))

	// 7. НАСТРОЙКА СЕРВЕРА С ТАЙМАУТАМИ
	srv := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	// 8. ЗАПУСК В ОТДЕЛЬНОЙ ГОРУТИНЕ
	go func() {
		fmt.Printf("Auth server running on port %s\n", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("Failed to start server: ", err)
		}
	}()

	// 9. GRACEFUL SHUTDOWN
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	fmt.Println("\nShutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server shutdown error: ", err)
	}
	fmt.Println("Server stopped")
}
