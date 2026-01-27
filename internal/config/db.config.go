package config

import (
	"fmt"
	"time"
)

type DBConfig struct {
	Host            string
	Port            int
	User            string
	Password        string
	DBName          string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
	ConnMaxIdleTime time.Duration
}

func NewDBConfig() (*DBConfig, error) {
	return &DBConfig{
		Host:            getEnv("DB_HOST", "postgres"),
		Port:            getEnvAsInt("DB_PORT", 5432),
		User:            getEnv("DB_USER", "user"),
		Password:        getEnv("DB_PASSWORD", "password"),
		DBName:          getEnv("DB_NAME", "dbname"),
		MaxOpenConns:    getEnvAsInt("DB_MAX_OPEN_CONNS", 25),
		MaxIdleConns:    getEnvAsInt("DB_MAX_IDLE_CONNS", 5),
		ConnMaxLifetime: time.Duration(getEnvAsInt("DB_CONN_MAX_LIFETIME_MIN", 30)) * time.Minute,
		ConnMaxIdleTime: time.Duration(getEnvAsInt("DB_CONN_MAX_IDLE_TIME_MIN", 5)) * time.Minute,
	}, nil
}

func (c *DBConfig) DSN() string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		c.Host, c.Port, c.User, c.Password, c.DBName,
	)
}
