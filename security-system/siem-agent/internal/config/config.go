package config

import (
	"encoding/json"
	"fmt"
	"os"
)

// Config — настройки агента. Живёт в пакете config, а не models,
// потому что это не доменная сущность (не то, чем "оперирует" бизнес-логика),
// а параметры поведения самого приложения.
type Config struct {
	LogPaths      []string `json:"log_paths"`
	ServerURL     string   `json:"server_url"`
	BatchSize     int      `json:"batch_size"`
	BatchInterval int      `json:"batch_interval"` // секунд
	RetryCount    int      `json:"retry_count"`
	RetryDelay    int      `json:"retry_delay"` // секунд
	FilterLevel   string   `json:"filter_level"`
}

// Load читает конфиг из JSON-файла. Если путь не указан — берёт дефолты.
func Load(path string) (*Config, error) {
	if path == "" {
		return DefaultConfig(), nil
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	cfg := DefaultConfig() // начинаем с дефолтов, JSON перезапишет только то, что в нём есть
// Unmarshal позволяет нам преобраазовать json в go выбирая только
// то что нам нужно тоесть мы избегаем попытки заалить скрипт нам в базу
	if err := json.Unmarshal(data, cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}
// Проверяем коректность конфига (код ниже через функцию)
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	return cfg, nil
}

func DefaultConfig() *Config {
	return &Config{
		LogPaths:      []string{"/var/log/syslog", "/var/log/auth.log"},
		ServerURL:     "http://localhost:8080/api/logs",
		BatchSize:     100,
		BatchInterval: 5,
		RetryCount:    3,
		RetryDelay:    5,
		FilterLevel:   "WARNING",
	}
}

// Validate проверяет, что значения имеют смысл — лучше упасть при старте,
// чем потом получить панику или бесконечный цикл где-то в середине работы агента.
func (c *Config) Validate() error {
	if len(c.LogPaths) == 0 {
		return fmt.Errorf("log_paths must not be empty")
	}
	if c.ServerURL == "" {
		return fmt.Errorf("server_url must not be empty")
	}
	if c.BatchSize <= 0 {
		return fmt.Errorf("batch_size must be positive, got %d", c.BatchSize)
	}
	if c.BatchInterval <= 0 {
		return fmt.Errorf("batch_interval must be positive, got %d", c.BatchInterval)
	}
	if c.RetryCount < 0 {
		return fmt.Errorf("retry_count must not be negative, got %d", c.RetryCount)
	}
	return nil
}