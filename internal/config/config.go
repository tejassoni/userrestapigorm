package config

import (
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

// Config holds all application configuration
type Config struct {
	App      AppConfig
	Server   ServerConfig
	API      APIConfig
	Database DatabaseConfig
	Logger   LoggerConfig
	JWT      JWTConfig
}

// AppConfig contains application metadata
type AppConfig struct {
	Name    string
	Version string
	Env     string // development, staging, production
}

// ServerConfig contains HTTP server settings
type ServerConfig struct {
	Host              string
	Port              string
	ReadTimeout       time.Duration
	WriteTimeout      time.Duration
	IdleTimeout       time.Duration
	ShutdownTimeout   time.Duration
	ReadHeaderTimeout time.Duration
}

// APIConfig contains the public API route prefix and version.
type APIConfig struct {
	Prefix  string
	Version string
}

// DatabaseConfig contains database connection settings
type DatabaseConfig struct {
	Driver          string // mysql, postgres
	Host            string
	Port            int
	Username        string
	Password        string
	Database        string
	MaxConnections  int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
}

// LoggerConfig contains logging settings
type LoggerConfig struct {
	Level  string // debug, info, warn, error
	Format string // json, text
}

// JWTConfig contains JWT settings
type JWTConfig struct {
	Secret string
	Expiry int // minutes
}

// Load loads configuration from environment variables
func Load() *Config {
	_ = godotenv.Load()

	return &Config{
		App: AppConfig{
			Name:    getenv("APP_NAME", "api"),
			Version: getenv("APP_VERSION", "1.0.0"),
			Env:     getenv("APP_ENV", "development"),
		},
		Server: ServerConfig{
			Host:              getenv("SERVER_HOST", ""),
			Port:              getenv("APP_PORT", "8080"),
			ReadHeaderTimeout: getenvDuration("SERVER_READ_HEADER_TIMEOUT", 5*time.Second),
			ReadTimeout:       getenvDuration("SERVER_READ_TIMEOUT", 10*time.Second),
			WriteTimeout:      getenvDuration("SERVER_WRITE_TIMEOUT", 10*time.Second),
			IdleTimeout:       getenvDuration("SERVER_IDLE_TIMEOUT", 60*time.Second),
			ShutdownTimeout:   getenvDuration("SERVER_SHUTDOWN_TIMEOUT", 30*time.Second),
		},
		API: APIConfig{
			Prefix:  getenv("API_PREFIX", "/api"),
			Version: getenv("API_VERSION", "v1"),
		},
		Database: DatabaseConfig{
			Driver:          getenv("DB_CONNECTION", "mysql"),
			Host:            getenv("DB_HOST", "localhost"),
			Port:            getenvi("DB_PORT", 3306),
			Username:        getenv("DB_USER", "root"),
			Password:        getenv("DB_PASSWORD", ""),
			Database:        getenv("DB_NAME", "api_db"),
			MaxConnections:  getenvi("DB_MAX_OPEN_CONNS", 10),
			MaxIdleConns:    getenvi("DB_MAX_IDLE_CONNS", 5),
			ConnMaxLifetime: getenvDuration("DB_CONN_MAX_LIFETIME", 60*time.Minute),
		},
		Logger: LoggerConfig{
			Level:  getenv("LOG_LEVEL", "info"),
			Format: getenv("LOG_FORMAT", "json"),
		},
		JWT: JWTConfig{
			Secret: getenv("JWT_SECRET", "your-secret-key-change-this"),
			Expiry: getenvi("JWT_EXPIRY", 15),
		},
	}
}

// Helper functions
func getenv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func getenvi(key string, defaultValue int) int {
	if value, exists := os.LookupEnv(key); exists {
		if intVal, err := strconv.Atoi(value); err == nil {
			return intVal
		}
	}
	return defaultValue
}

func getenvDuration(key string, defaultValue time.Duration) time.Duration {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}

	duration, err := time.ParseDuration(value)
	if err != nil {
		return defaultValue
	}

	return duration
}
