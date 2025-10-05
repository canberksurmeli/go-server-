package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/go-playground/validator/v10"
)

type BaseConfig struct {
	LogLevel string `validate:"required,oneof=debug info warn error"`
}

type HttpConfig struct {
	Port int `validate:"required,min=1,max=65535"`
}

type DatabaseConfig struct {
	Host     string `validate:"required,hostname|ip"`
	Port     int    `validate:"required,min=1,max=65535"`
	User     string `validate:"required"`
	Password string `validate:"required"`
	Name     string `validate:"required"`
}

type Config struct {
	Http HttpConfig
	Base BaseConfig
	Database DatabaseConfig
}

var instance *Config

func Load() *Config {
	port := 8080
	if p := os.Getenv("PORT"); p != "" {
		if parsed, err := strconv.Atoi(p); err == nil {
			port = parsed
		}
	}

	logLevel := os.Getenv("LOG_LEVEL")
	if logLevel == "" {
		logLevel = "info"
	}

	env := os.Getenv("ENV")
	if env == "" {
		env = "development"
	}

	config := &Config{
		Http: HttpConfig{
			Port: port,
		},
		Base: BaseConfig{
			LogLevel: logLevel,
		},
		Database: DatabaseConfig{
			Host:     os.Getenv("DB_HOST"),
			Port:     getEnvAsInt("DB_PORT", 5432),
			User:     getEnvOrDefault("DB_USER", "postgres"),
			Password: getEnvOrDefault("DB_PASSWORD", "password"),
			Name:     getEnvOrDefault("DB_NAME", "message_provider"),
		},
	}

	validate := validator.New(validator.WithRequiredStructEnabled())

	err := validate.Struct(config)
	if err != nil {
		fmt.Println("Configuration validation error:", err)
	}

	return config
}

func Get() *Config {
	if instance == nil {
		instance = Load()
	}
	return instance
}

func getEnvAsInt(name string, defaultVal int) int {
	if valStr := os.Getenv(name); valStr != "" {
		if val, err := strconv.Atoi(valStr); err == nil {
			return val
		}
	}
	return defaultVal
}

func getEnvOrDefault(name, defaultVal string) string {
	if val := os.Getenv(name); val != "" {
		return val
	}
	return defaultVal
}
