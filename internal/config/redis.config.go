package config

import "fmt"

type RedisConfig struct {
	Host string
	Port int
}

func NewRedisConfig() (*RedisConfig, error) {
	return &RedisConfig{
		Host: getEnv("REDIS_HOST", "redis"),
		Port: getEnvAsInt("REDIS_PORT", 6379),
	}, nil
}

func (c *RedisConfig) Addr() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}
