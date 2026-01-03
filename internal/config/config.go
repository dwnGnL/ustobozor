package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

// Config holds application configuration loaded from environment variables.
type Config struct {
	App    AppConfig
	DB     DBConfig
	JWT    JWTConfig
	SMS    SMSConfig
	Upload UploadConfig
}

type AppConfig struct {
	Env      string
	Port     string
	LogLevel string
}

type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	SSLMode  string
}

type JWTConfig struct {
	AccessSecret  string
	RefreshSecret string
	AccessTTL     time.Duration
	RefreshTTL    time.Duration
}

type SMSConfig struct {
	Provider string
	APIKey   string
	Sender   string
}

type UploadConfig struct {
	Dir          string
	MaxSizeBytes int64
}

func Load() (*Config, error) {
	_ = godotenv.Load()

	cfg := &Config{
		App: AppConfig{
			Env:      getEnv("APP_ENV", "local"),
			Port:     getEnv("APP_PORT", "8080"),
			LogLevel: getEnv("APP_LOG_LEVEL", "info"),
		},
		DB: DBConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", "postgres"),
			Name:     getEnv("DB_NAME", "bozor"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
		},
		JWT: JWTConfig{
			AccessSecret:  getEnv("JWT_ACCESS_SECRET", "dev_access_secret"),
			RefreshSecret: getEnv("JWT_REFRESH_SECRET", "dev_refresh_secret"),
			AccessTTL:     getEnvDuration("JWT_ACCESS_TTL", 24*time.Hour),
			RefreshTTL:    getEnvDuration("JWT_REFRESH_TTL", 30*24*time.Hour),
		},
		SMS: SMSConfig{
			Provider: getEnv("SMS_PROVIDER", "mock"),
			APIKey:   getEnv("SMS_API_KEY", ""),
			Sender:   getEnv("SMS_SENDER", "BOZOR"),
		},
		Upload: UploadConfig{
			Dir:          getEnv("UPLOAD_DIR", "./uploads"),
			MaxSizeBytes: getEnvInt64("UPLOAD_MAX_MB", 25) * 1024 * 1024,
		},
	}

	if cfg.JWT.AccessSecret == cfg.JWT.RefreshSecret {
		return nil, fmt.Errorf("jwt access and refresh secrets must differ")
	}

	return cfg, nil
}

func (c DBConfig) DSN() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		c.User,
		c.Password,
		c.Host,
		c.Port,
		c.Name,
		c.SSLMode,
	)
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}

func getEnvDuration(key string, fallback time.Duration) time.Duration {
	val := os.Getenv(key)
	if val == "" {
		return fallback
	}
	parsed, err := time.ParseDuration(val)
	if err != nil {
		return fallback
	}
	return parsed
}

func getEnvInt64(key string, fallback int64) int64 {
	val := os.Getenv(key)
	if val == "" {
		return fallback
	}
	parsed, err := strconv.ParseInt(val, 10, 64)
	if err != nil {
		return fallback
	}
	return parsed
}
