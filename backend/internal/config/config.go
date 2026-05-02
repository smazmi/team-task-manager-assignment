package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

type Config struct {
	AppPort     string
	GinMode     string
	DatabaseURL string
	DBHost      string
	DBPort      string
	DBUser      string
	DBPassword  string
	DBName      string
	DBSSLMode   string
	JWTSecret   string
	JWTTTL      time.Duration
	FrontendURL string
	AutoMigrate bool
}

func Load() (Config, error) {
	jwtTTLHours, err := strconv.Atoi(getEnv("JWT_TTL_HOURS", "24"))
	if err != nil || jwtTTLHours < 1 {
		return Config{}, fmt.Errorf("JWT_TTL_HOURS must be a positive integer")
	}

	autoMigrate, err := strconv.ParseBool(getEnv("AUTO_MIGRATE", "true"))
	if err != nil {
		return Config{}, fmt.Errorf("AUTO_MIGRATE must be true or false")
	}

	return Config{
		AppPort:     getEnv("APP_PORT", "8080"),
		GinMode:     getEnv("GIN_MODE", "debug"),
		DatabaseURL: os.Getenv("DATABASE_URL"),
		DBHost:      getEnv("DB_HOST", "localhost"),
		DBPort:      getEnv("DB_PORT", "5432"),
		DBUser:      getEnv("DB_USER", "postgres"),
		DBPassword:  getEnv("DB_PASSWORD", "postgres"),
		DBName:      getEnv("DB_NAME", "team_task_manager"),
		DBSSLMode:   getEnv("DB_SSL_MODE", "disable"),
		JWTSecret:   getEnv("JWT_SECRET", "change-me-in-production"),
		JWTTTL:      time.Duration(jwtTTLHours) * time.Hour,
		FrontendURL: getEnv("FRONTEND_URL", "http://localhost:5173"),
		AutoMigrate: autoMigrate,
	}, nil
}

func (c Config) DatabaseDSN() string {
	if c.DatabaseURL != "" {
		return c.DatabaseURL
	}

	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.DBHost,
		c.DBPort,
		c.DBUser,
		c.DBPassword,
		c.DBName,
		c.DBSSLMode,
	)
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok && value != "" {
		return value
	}

	return fallback
}
