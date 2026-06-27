// Package config loads and validates all application configuration
// from environment variables following the 12-factor app methodology.
package config

import (
	"fmt"
	"os"
)

// Config holds all application configuration.
type Config struct {
	App      AppConfig
	Server   ServerConfig
	Database DatabaseConfig
}

// AppConfig holds general application settings.
type AppConfig struct {
	// APP_ENV: application environment (development, staging, production)
	Env string
	// APP_NAME: application name used in logs
	Name string
}

// ServerConfig holds HTTP server settings.
type ServerConfig struct {
	// SERVER_PORT: port the HTTP server listens on (default: 8080)
	Port string
	// SERVER_READ_TIMEOUT: max time to read request (default: 15s)
	ReadTimeout string
	// SERVER_WRITE_TIMEOUT: max time to write response (default: 15s)
	WriteTimeout string
}

// DatabaseConfig holds PostgreSQL connection settings.
type DatabaseConfig struct {
	// DB_HOST: postgres host (default: localhost)
	Host string
	// DB_PORT: postgres port (default: 5432)
	Port string
	// DB_NAME: database name (required)
	Name string
	// DB_USER: database user (required)
	User string
	// DB_PASSWORD: database password (required)
	Password string
	// DB_SSLMODE: ssl mode (default: disable)
	SSLMode string
	// DB_MAX_OPEN_CONNS: max open connections (default: 25)
	MaxOpenConns string
	// DB_MAX_IDLE_CONNS: max idle connections (default: 25)
	MaxIdleConns string
	// DB_CONN_MAX_LIFETIME: connection max lifetime (default: 5m)
	ConnMaxLifetime string
}

// DSN returns the PostgreSQL data source name.
func (d DatabaseConfig) DSN() string {
	return fmt.Sprintf(
		"host=%s port=%s dbname=%s user=%s password=%s sslmode=%s",
		d.Host, d.Port, d.Name, d.User, d.Password, d.SSLMode,
	)
}

// Load reads configuration from environment variables.
// Returns an error if any required variable is missing.
func Load() (*Config, error) {
	dbName, err := requireEnv("DB_NAME")
	if err != nil {
		return nil, err
	}
	dbUser, err := requireEnv("DB_USER")
	if err != nil {
		return nil, err
	}
	dbPassword, err := requireEnv("DB_PASSWORD")
	if err != nil {
		return nil, err
	}

	return &Config{
		App: AppConfig{
			Env:  getEnv("APP_ENV", "development"),
			Name: getEnv("APP_NAME", "task-api"),
		},
		Server: ServerConfig{
			Port:         getEnv("SERVER_PORT", "8080"),
			ReadTimeout:  getEnv("SERVER_READ_TIMEOUT", "15s"),
			WriteTimeout: getEnv("SERVER_WRITE_TIMEOUT", "15s"),
		},
		Database: DatabaseConfig{
			Host:            getEnv("DB_HOST", "localhost"),
			Port:            getEnv("DB_PORT", "5432"),
			Name:            dbName,
			User:            dbUser,
			Password:        dbPassword,
			SSLMode:         getEnv("DB_SSLMODE", "disable"),
			MaxOpenConns:    getEnv("DB_MAX_OPEN_CONNS", "25"),
			MaxIdleConns:    getEnv("DB_MAX_IDLE_CONNS", "25"),
			ConnMaxLifetime: getEnv("DB_CONN_MAX_LIFETIME", "5m"),
		},
	}, nil
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func requireEnv(key string) (string, error) {
	v := os.Getenv(key)
	if v == "" {
		return "", fmt.Errorf("required environment variable %q is not set", key)
	}
	return v, nil
}
