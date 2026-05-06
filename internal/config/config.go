package config

import "os"

type Config struct {
	AppName     string
	AppPort     string
	AppEnv      string
	DatabaseURL string
}

func Load() Config {
	return Config{
		AppName:     getEnv("APP_NAME", "Gogal Studio"),
		AppPort:     getEnv("APP_PORT", "8080"),
		AppEnv:      getEnv("APP_ENV", "development"),
		DatabaseURL: getEnv("DATABASE_URL", "postgres://door:door@localhost:5432/door_app?sslmode=disable"),
	}
}

func getEnv(key string, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	return value
}
