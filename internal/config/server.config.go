package config

type ServerConfig struct {
	Environment string
	LogLevel    string
	Port        string
}

func NewServerConfig() (*ServerConfig, error) {
	return &ServerConfig{
		Environment: getEnv("APP_ENV", "development"),
		Port:        getEnv("APP_PORT", "8000"),
		LogLevel:    getEnv("LOG_LEVEL", "info"),
	}, nil
}
