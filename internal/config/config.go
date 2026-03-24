package config

import "os"

type Config struct {
	HTTPAddr string
	DBPath   string
}

func New() Config {
	return Config{
		HTTPAddr: getEnv("HTTP_ADDR", ":8080"),
		DBPath:   getEnv("SQLITE_PATH", "./data/demo.db"),
	}
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}

	return fallback
}
