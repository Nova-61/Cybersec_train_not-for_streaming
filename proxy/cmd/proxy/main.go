package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"proxy/internal/config"
	"proxy/internal/handler"
	"proxy/internal/middleware"
)

func loadDotEnv(path string) {
	f, err := os.Open(path)
	if err != nil {
		return
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}
		key := strings.TrimSpace(parts[0])
		val := strings.TrimSpace(parts[1])
		if _, exists := os.LookupEnv(key); !exists {
			os.Setenv(key, val)
		}
	}
}

func getEnv(key, def string) string {
	if v, ok := os.LookupEnv(key); ok && v != "" {
		return v
	}
	return def
}

func getEnvInt(key string, def int) int {
	if v, ok := os.LookupEnv(key); ok && v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			return n
		}
	}
	return def
}

func main() {
	loadDotEnv(".env")
	listenAddr := flag.String("listen", getEnv("LISTEN_ADDR", ":8888"), "адрес для прокси")
	targetURL := flag.String("target", getEnv("TARGET_URL", "http://localhost:8080"), "целевой сервер")
	rateLimit := flag.Int("rate-limit", getEnvInt("RATE_LIMIT", 10), "максимум запросов в минуту")
	flag.Parse()

	cfg := config.New()
	cfg.ListenAddr = *listenAddr
	cfg.TargetURL = *targetURL

	proxyHandler := handler.NewProxyHandler(cfg.TargetURL)

	limiter := middleware.NewRateLimiter(*rateLimit, time.Minute)
	handlerChain := middleware.RateLimit(limiter)(
		middleware.Security(
			middleware.Logger(proxyHandler),
		),
	)

	// Создаём сервер и явно назначаем handlerChain
	srv := &http.Server{
		Addr:         cfg.ListenAddr,
		Handler:      handlerChain,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	go func() {
		fmt.Printf("🚀 Прокси запущен на %s\n", cfg.ListenAddr)
		fmt.Printf("🎯 Перенаправляет на %s\n", cfg.TargetURL)
		fmt.Printf("⏱️  Rate Limit: %d запросов в минуту\n", *rateLimit)
		fmt.Println("📝 Нажми Ctrl+C для выхода")
		fmt.Println()

		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("Ошибка запуска сервера:", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	fmt.Println("\n🛑 Получен сигнал завершения. Останавливаем сервер...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Ошибка при завершении:", err)
	}

	fmt.Println("✅ Сервер корректно завершён")
}
