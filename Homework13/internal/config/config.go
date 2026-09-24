package config

import (
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"
)

// Config хранит настройки, загруженные из .env
type Config struct {
	Login     string
	Password  string
	JWTSecret []byte
	JWTTTL    time.Duration
}

// Load читает .env и переменные окружения, возвращает Config.
func Load() (*Config, error) {
	_ = godotenv.Load()

	login := os.Getenv("LOGIN")
	if login == "" {
		return nil, fmt.Errorf("переменная LOGIN не задана")
	}

	password := os.Getenv("PASSWORD")
	if password == "" {
		return nil, fmt.Errorf("переменная PASSWORD не задана")
	}

	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return nil, fmt.Errorf("переменная JWT_SECRET не задана")
	}

	ttlStr := os.Getenv("JWT_TTL")
	if ttlStr == "" {
		return nil, fmt.Errorf("переменная JWT_TTL не задана")
	}
	ttl, err := time.ParseDuration(ttlStr)
	if err != nil {
		return nil, fmt.Errorf("неверный формат JWT_TTL: %w", err)
	}

	return &Config{
		Login:     login,
		Password:  password,
		JWTSecret: []byte(secret),
		JWTTTL:    ttl,
	}, nil
}
