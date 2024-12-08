package config

import (
	"os"
	"time"
)

type Config struct {
	Storage
	Server
	Logger
}

type Storage struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

type Server struct {
	Address     string
	Timeout     time.Duration
	IdleTimeout time.Duration
}

type Logger struct {
	Env string
}

func MustLoad() *Config {
	return &Config{
		Storage: Storage{
			Host:     getEnv("DB_HOST", ""),
			Port:     getEnv("DB_PORT", ""),
			User:     getEnv("DB_USER", ""),
			Password: getEnv("DB_PASSWORD", ""),
			DBName:   getEnv("DB_NAME", ""),
			SSLMode:  getEnv("DB_SSLMODE", ""),
		},
		Server: Server{
			Address:     getEnv("address", ""),
			Timeout:     getEnvAsTime("timeout", "4s"),
			IdleTimeout: getEnvAsTime("idle_timeout", "60s"),
		},
		Logger: Logger{
			Env: getEnv("env", "debug"),
		},
	}
}

func getEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return defaultVal
}

func getEnvAsTime(key string, defaultVal string) time.Duration {
	if value, exists := os.LookupEnv(key); exists {
		duration, err := time.ParseDuration(value)
		if err == nil {
			return duration
		}
	}

	duration, _ := time.ParseDuration(defaultVal)
	return duration
}
