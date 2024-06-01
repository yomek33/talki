package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

// Config holds the application configuration
type Config struct {
	TiDBUser     string
	TiDBPassword string
	TiDBHost     string
	TiDBPort     string
	TiDBDBName   string
	UseSSL       string
	Port         string
	GeminiAPIKey string
}

// LoadConfig loads configuration from environment variables
func LoadConfig() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		return nil, fmt.Errorf("error loading .env file: %w", err)
	}

	cfg := &Config{
		TiDBUser:     os.Getenv("TIDB_USER"),
		TiDBPassword: os.Getenv("TIDB_PASSWORD"),
		TiDBHost:     os.Getenv("TIDB_HOST"),
		TiDBPort:     os.Getenv("TIDB_PORT"),
		TiDBDBName:   os.Getenv("TIDB_DB_NAME"),
		UseSSL:       os.Getenv("USE_SSL"),
		Port:         os.Getenv("PORT"),
		GeminiAPIKey: os.Getenv("GEMINI_API_KEY"),
	}

	if cfg.TiDBUser == "" || cfg.TiDBPassword == "" || cfg.TiDBHost == "" || cfg.TiDBPort == "" || cfg.TiDBDBName == "" || cfg.Port == "" || cfg.UseSSL == "" || cfg.GeminiAPIKey == ""{
		return nil, fmt.Errorf("one or more required environment variables are missing")
	}

	return cfg, nil
}
