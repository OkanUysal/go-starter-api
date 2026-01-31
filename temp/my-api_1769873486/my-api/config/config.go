package config

import "os"

type Config struct {
	AppName string
	Port    string
	DatabaseURL string
	JWTSecret string
	LogLevel string
}

func Load() *Config {
	return &Config{
		AppName: getEnv("APP_NAME", "my-api"),
		Port: getEnv("PORT", "8080"),
		DatabaseURL: getEnv("DATABASE_URL", ""),
		JWTSecret: getEnv("JWT_SECRET", ""),
		LogLevel: getEnv("LOG_LEVEL", "info"),
	}
}

func getEnv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
