package config

import (
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

// Config holds all application configuration
type Config struct {
	App        AppConfig
	Server     ServerConfig
	API        APIConfig
	Database   DatabaseConfig
	Redis      RedisConfig
	Memcached  MemcachedConfig
	SMTP       SMTPConfig
	Logger     LoggerConfig
	Upload     UploadConfig
	Pagination PaginationConfig
	RateLimit  RateLimitConfig
	JWT        JWTConfig
	AWS        AWSConfig
	CORS       CORSConfig
	Features   FeatureConfig
}

// AppConfig contains application metadata
type AppConfig struct {
	Name           string
	Version        string
	Env            string
	Debug          bool
	Timezone       string
	URL            string
	Locale         string
	FallbackLocale string
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
	MaxHeaderBytes    int
	MaxBodySize       int64
	DisableKeepAlives bool
}

// APIConfig contains API route settings
type APIConfig struct {
	Prefix  string
	Version string
	Auth    string
}

// DatabaseConfig contains database connection settings
type DatabaseConfig struct {
	Driver          string
	Host            string
	Port            int
	Username        string
	Password        string
	Database        string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
	ConnMaxIdleTime time.Duration
}

// RedisConfig contains Redis connection settings
type RedisConfig struct {
	Host     string
	Port     int
	Password string
	DB       int
}

// MemcachedConfig contains Memcached connection settings
type MemcachedConfig struct {
	Host string
	Port int
}

// SMTPConfig contains email settings
type SMTPConfig struct {
	Host      string
	Port      int
	Username  string
	Password  string
	FromName  string
	FromEmail string
}

// LoggerConfig contains logging settings
type LoggerConfig struct {
	Level  string
	Format string
}

// UploadConfig contains file upload settings
type UploadConfig struct {
	Directory string
	MaxSize   int64
}

// PaginationConfig contains pagination defaults
type PaginationConfig struct {
	DefaultSize int
	MaxSize     int
}

// RateLimitConfig contains rate limiting settings
type RateLimitConfig struct {
	Requests int
	Duration time.Duration
}

// JWTConfig contains JWT settings
type JWTConfig struct {
	Secret     string
	Expiration time.Duration
	Issuer     string
}

// AWSConfig contains AWS service settings
type AWSConfig struct {
	AccessKeyID     string
	SecretAccessKey string
	DefaultRegion   string
	Bucket          string
	UsePathStyle    bool
}

// CORSConfig contains CORS settings
type CORSConfig struct {
	AllowedOrigins []string
	AllowedMethods []string
	AllowedHeaders []string
}

// FeatureConfig contains feature flags
type FeatureConfig struct {
	Cache     bool
	Metrics   bool
	Tracing   bool
	RateLimit bool
}

// Load loads configuration from environment variables
func Load() *Config {
	_ = godotenv.Load()

	return &Config{
		App: AppConfig{
			Name:           getEnv("APP_NAME", "api"),
			Version:        getEnv("APP_VERSION", "1.0.0"),
			Env:            getEnv("APP_ENV", "development"),
			Debug:          getEnvBool("APP_DEBUG", true),
			Timezone:       getEnv("APP_TIMEZONE", "UTC"),
			URL:            getEnv("APP_URL", "http://localhost:8080"),
			Locale:         getEnv("APP_LOCALE", "en"),
			FallbackLocale: getEnv("APP_FALLBACK_LOCALE", "en"),
		},
		Server: ServerConfig{
			Host:              getEnv("SERVER_HOST", "127.0.0.1"),
			Port:              getEnv("SERVER_PORT", "8080"),
			ReadHeaderTimeout: getEnvDuration("SERVER_READ_HEADER_TIMEOUT", 5*time.Second),
			ReadTimeout:       getEnvDuration("SERVER_READ_TIMEOUT", 10*time.Second),
			WriteTimeout:      getEnvDuration("SERVER_WRITE_TIMEOUT", 10*time.Second),
			IdleTimeout:       getEnvDuration("SERVER_IDLE_TIMEOUT", 60*time.Second),
			ShutdownTimeout:   getEnvDuration("SERVER_SHUTDOWN_TIMEOUT", 30*time.Second),
			MaxHeaderBytes:    getEnvInt("SERVER_MAX_HEADER_BYTES", 1048576),
			MaxBodySize:       getEnvInt64("SERVER_MAX_REQUEST_BODY_SIZE", 10485760),
			DisableKeepAlives: getEnvBool("SERVER_DISABLE_KEEP_ALIVES", false),
		},
		API: APIConfig{
			Prefix:  getEnv("API_PREFIX", "/api"),
			Version: getEnv("API_VERSION", "v1"),
			Auth:    getEnv("API_AUTH", "/auth"),
		},
		Database: DatabaseConfig{
			Driver:          getEnv("DB_DRIVER", "mysql"),
			Host:            getEnv("DB_HOST", "127.0.0.1"),
			Port:            getEnvInt("DB_PORT", 3306),
			Username:        getEnv("DB_USER", "root"),
			Password:        getEnv("DB_PASSWORD", ""),
			Database:        getEnv("DB_NAME", "api_db"),
			MaxOpenConns:    getEnvInt("DB_MAX_OPEN_CONNS", 10),
			MaxIdleConns:    getEnvInt("DB_MAX_IDLE_CONNS", 10),
			ConnMaxLifetime: getEnvDuration("DB_CONN_MAX_LIFETIME", 5*time.Minute),
			ConnMaxIdleTime: getEnvDuration("DB_CONN_MAX_IDLE_TIME", 2*time.Minute),
		},
		Redis: RedisConfig{
			Host:     getEnv("REDIS_HOST", "127.0.0.1"),
			Port:     getEnvInt("REDIS_PORT", 6379),
			Password: getEnv("REDIS_PASSWORD", ""),
			DB:       getEnvInt("REDIS_DB", 0),
		},
		Memcached: MemcachedConfig{
			Host: getEnv("MEMCACHED_HOST", "127.0.0.1"),
			Port: getEnvInt("MEMCACHED_PORT", 11211),
		},
		SMTP: SMTPConfig{
			Host:      getEnv("SMTP_HOST", ""),
			Port:      getEnvInt("SMTP_PORT", 587),
			Username:  getEnv("SMTP_USERNAME", ""),
			Password:  getEnv("SMTP_PASSWORD", ""),
			FromName:  getEnv("SMTP_FROM_NAME", ""),
			FromEmail: getEnv("SMTP_FROM_EMAIL", ""),
		},
		Logger: LoggerConfig{
			Level:  getEnv("LOG_LEVEL", "info"),
			Format: getEnv("LOG_FORMAT", "json"),
		},
		Upload: UploadConfig{
			Directory: getEnv("UPLOAD_DIR", "./storage/uploads"),
			MaxSize:   getEnvInt64("UPLOAD_MAX_SIZE", 10485760),
		},
		Pagination: PaginationConfig{
			DefaultSize: getEnvInt("PAGINATION_DEFAULT_SIZE", 10),
			MaxSize:     getEnvInt("PAGINATION_MAX_SIZE", 100),
		},
		RateLimit: RateLimitConfig{
			Requests: getEnvInt("RATE_LIMIT_REQUESTS", 100),
			Duration: getEnvDuration("RATE_LIMIT_DURATION", 1*time.Minute),
		},
		JWT: JWTConfig{
			Secret:     getEnv("JWT_SECRET", "your-secret-key-change-this"),
			Expiration: getEnvDuration("JWT_EXPIRATION", 24*time.Hour),
			Issuer:     getEnv("JWT_ISSUER", "UserRestAPIGo"),
		},
		AWS: AWSConfig{
			AccessKeyID:     getEnv("AWS_ACCESS_KEY_ID", ""),
			SecretAccessKey: getEnv("AWS_SECRET_ACCESS_KEY", ""),
			DefaultRegion:   getEnv("AWS_DEFAULT_REGION", "us-east-1"),
			Bucket:          getEnv("AWS_BUCKET", ""),
			UsePathStyle:    getEnvBool("AWS_USE_PATH_STYLE_ENDPOINT", false),
		},
		CORS: CORSConfig{
			AllowedOrigins: getEnvSlice("CORS_ALLOWED_ORIGINS", []string{"*"}),
			AllowedMethods: getEnvSlice("CORS_ALLOWED_METHODS", []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"}),
			AllowedHeaders: getEnvSlice("CORS_ALLOWED_HEADERS", []string{"Content-Type", "Authorization"}),
		},
		Features: FeatureConfig{
			Cache:     getEnvBool("FEATURE_CACHE", true),
			Metrics:   getEnvBool("FEATURE_METRICS", false),
			Tracing:   getEnvBool("FEATURE_TRACING", false),
			RateLimit: getEnvBool("FEATURE_RATE_LIMIT", true),
		},
	}
}

// Helper functions
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists && value != "" {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value, exists := os.LookupEnv(key); exists {
		if intVal, err := strconv.Atoi(value); err == nil {
			return intVal
		}
	}
	return defaultValue
}

func getEnvInt64(key string, defaultValue int64) int64 {
	if value, exists := os.LookupEnv(key); exists {
		if intVal, err := strconv.ParseInt(value, 10, 64); err == nil {
			return intVal
		}
	}
	return defaultValue
}

func getEnvBool(key string, defaultValue bool) bool {
	if value, exists := os.LookupEnv(key); exists {
		if boolVal, err := strconv.ParseBool(value); err == nil {
			return boolVal
		}
	}
	return defaultValue
}

func getEnvDuration(key string, defaultValue time.Duration) time.Duration {
	if value, exists := os.LookupEnv(key); exists && value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return defaultValue
}

func getEnvSlice(key string, defaultValue []string) []string {
	if value, exists := os.LookupEnv(key); exists && value != "" {
		parts := strings.Split(value, ",")
		result := make([]string, 0, len(parts))
		for _, part := range parts {
			if trimmed := strings.TrimSpace(part); trimmed != "" {
				result = append(result, trimmed)
			}
		}
		if len(result) > 0 {
			return result
		}
	}
	return defaultValue
}
