package config

import "os"

type Config struct {
	Port      string
	JWTSecret string
	DBPath    string
}

func Load() *Config {
	return &Config{
		Port:      getEnv("PORT", "8080"),
		JWTSecret: getEnv("JWT_SECRET", "smart-cashier-secret-key-2026"),
		DBPath:    getEnv("DB_PATH", "smart_cashier.db"),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
