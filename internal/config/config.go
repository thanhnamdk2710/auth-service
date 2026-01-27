package config

import (
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	Server *ServerConfig
	DB     *DBConfig
	Redis  *RedisConfig
}

func NewConfig() (*Config, error) {
	serverConfig, err := NewServerConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to load server config: %w", err)
	}

	dbConfig, err := NewDBConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to load database config: %w", err)
	}

	redisConfig, err := NewRedisConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to load redis config: %w", err)
	}

	return &Config{
		DB:     dbConfig,
		Redis:  redisConfig,
		Server: serverConfig,
	}, nil
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
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
