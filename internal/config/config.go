package config

import (
	"os"
	"strconv"
)

type Config struct {
	DBHost        string
	DBPort        string
	DBUser        string
	DBPassword    string
	DBName        string
	ServerPort    string
	AuthSecret    string
	PasswordSecret string
	AuthTTLMinutes int
	Debug         bool
}

func Load() *Config {
	return &Config{
		DBHost:         getEnv("SQL_HOST", "localhost"),
		DBPort:         getEnv("SQL_PORT", "5432"),
		DBUser:         getEnv("SQL_USER", "surl_u"),
		DBPassword:     getEnv("SQL_PASSWORD", "surl_p"),
		DBName:         getEnv("SQL_DATABASE", "surl_d"),
		ServerPort:     getEnv("SERVER_PORT", "8000"),
		AuthSecret:     getEnv("AUTH_SECRET", "TODO replace with a random secure string 123"),
		PasswordSecret: getEnv("PASSWORD_SECRET", "TODO here replace with a random secure string or place in db 456"),
		AuthTTLMinutes: getEnvInt("AUTH_TTL_MINUTES", 12000),
		Debug:          getEnvBool("DEBUG", false),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	if v := os.Getenv(key); v != "" {
		if i, err := strconv.Atoi(v); err == nil {
			return i
		}
	}
	return fallback
}

func getEnvBool(key string, fallback bool) bool {
	if v := os.Getenv(key); v != "" {
		if b, err := strconv.ParseBool(v); err == nil {
			return b
		}
	}
	return fallback
}
