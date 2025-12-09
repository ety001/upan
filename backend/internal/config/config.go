package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Port          string
	FileMaxSize   int64 // in bytes
	FileExpireTime int  // in hours
	StoragePath   string
	DBPath        string
}

var AppConfig *Config

func Load() error {
	// Load .env file if it exists
	_ = godotenv.Load()

	AppConfig = &Config{
		Port:          getEnv("PORT", "8080"),
		FileMaxSize:   int64(getEnvAsInt("FILE_MAX_SIZE", 30) * 1024 * 1024), // Convert MB to bytes
		FileExpireTime: getEnvAsInt("FILE_EXPIRE_TIME", 6),
		StoragePath:   getEnv("STORAGE_PATH", "./storage/files"),
		DBPath:        getEnv("DB_PATH", "./storage/database.db"),
	}

	// Ensure storage directory exists
	if err := os.MkdirAll(AppConfig.StoragePath, 0755); err != nil {
		return err
	}

	return nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return defaultValue
	}
	value, err := strconv.Atoi(valueStr)
	if err != nil {
		return defaultValue
	}
	return value
}

