package config

import (
	"os"
)

type Config struct {
	Port      string
	DBPath    string
	SecretKey string
}

func Load() *Config {
	return &Config{
		Port:      getEnv("PORT", "8080"),
		DBPath:    getEnv("DB_PATH", "./data/repair.db"),
		SecretKey: getEnv("SECRET_KEY", "super-secret-key-change-in-prod"),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
