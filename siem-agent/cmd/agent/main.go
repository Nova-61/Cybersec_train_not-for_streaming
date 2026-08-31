package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"siem-agent/internal/config"
	"siem-agent/internal/filter"
	"siem-agent/internal/models"
	"siem-agent/internal/parser"
	"siem-agent/internal/reader"
	"siem-agent/internal/sender"
)

func main() {
	configPath := flag.String("config", "", "путь к конфиг-файлу")
	flag.Parse()

	cfg, err := config.Load(*configPath)
	if err != nil {
		log.Fatal("Failed to load config: ", err)
	}

	fmt.Printf("SIEM Agent starting...\n")
	fmt.Printf("Log paths: %v\n", cfg.LogPaths)
	fmt.Printf("Server URL: %s\n", cfg.ServerURL)
	fmt.Printf("Filter level: %s\n", cfg.FilterLevel)
	fmt.Println()

	logReader, err := reader.NewLogReader(cfg.LogPaths)
	if err != nil {
		log.Fatal("Failed to create reader: ", err)
	}

	logParser := parser.NewParser()
	logFilter := filter.NewFilter(cfg.FilterLevel)
	logSender := sender.NewSender(cfg.ServerURL, cfg.RetryCount, cfg.RetryDelay)

	logReader.Start()

	logChan := logReader.Output()
	errChan := logReader.Errors()

	// Сигналы ОС слушаем через отдельный канал, подписываемся ОДИН раз при старте -
	// не внутри select, а до него (в отличие от исходной версии)
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	batch := make([]models.LogEntry, 0, cfg.BatchSize)
	ticker := time.NewTicker(time.Duration(cfg.BatchInterval) * time.Second)
	defer ticker.Stop()

	// flush отправляет текущий батч и очищает его. Вынесено в отдельную функцию,
	// потому что вызывается из трёх разных мест select'а — не дублируем логику.
	flush := func(reason string) {
		if len(batch) == 0 {
			return
		}
		if err := logSender.SendBatch(batch); err != nil {
			fmt.Printf("Failed to send batch (%s): %v\n", reason, err)
		} else {
			fmt.Printf("Sent %d logs (%s)\n", len(batch), reason)
		}
		batch = make([]models.LogEntry, 0, cfg.BatchSize)
	}

	fmt.Println("Agent is running. Press Ctrl+C to stop.")

	for {
		select {
		case entry := <-logChan:
			parsed := logParser.Parse(entry.Message, entry.Source)
			if !logFilter.ShouldSend(parsed) {
				continue
			}
			batch = append(batch, parsed)
			if len(batch) >= cfg.BatchSize {
				flush("size limit")
			}

		case <-ticker.C:
			flush("timer")

		case err := <-errChan:
			fmt.Printf("Reader error: %v\n", err)

		case <-stop:
			fmt.Println("\nShutting down...")
			flush("shutdown") // не теряем то, что накопилось в батче на момент остановки
			logReader.Close()
			return
		}
	}
}
