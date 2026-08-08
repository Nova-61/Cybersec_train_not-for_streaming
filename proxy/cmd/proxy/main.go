package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"proxy/internal/config"
	"proxy/internal/handler"
	"proxy/internal/middleware"
)

func main() {
	listenAddr := flag.String("listen", ":8888", "адрес для прокси")
	targetURL := flag.String("target", "http://localhost:8080", "целевой сервер")
	rateLimit := flag.Int("rate-limit", 10, "максимум запросов в минуту")
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

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		handlerChain.ServeHTTP(w, r)
	})

	srv := &http.Server{
		Addr:         cfg.ListenAddr,
		Handler:      nil,
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
