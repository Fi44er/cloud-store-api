// Пакет config содержит конфигурацию приложения
// Загрузка происходит из переменных окружения
package config

import (
	"os"
	"strconv"
)

// Config — основная структура конфигурации приложения
type Config struct {
	DatabaseURL   string
	JWTSecret     string
	SessionSecret string
	UploadDir     string
	MaxQuotaBytes int64
	Port          string
	Env           string
}

// Load загружает конфигурацию из переменных окружения
func Load() *Config {
	maxQuota, err := strconv.ParseInt(getEnv("MAX_QUOTA_BYTES", "1073741824"), 10, 64)
	if err != nil {
		maxQuota = 1073741824 // 1 ГБ по умолчанию
	}

	return &Config{
		DatabaseURL:   getEnv("DATABASE_URL", "postgres://gloude:gloude_secret@localhost:5432/gloude_store?sslmode=disable"),
		JWTSecret:     getEnv("JWT_SECRET", "super_secret_jwt_key"),
		SessionSecret: getEnv("SESSION_SECRET", "super_secret_session"),
		UploadDir:     getEnv("UPLOAD_DIR", "./uploads"),
		MaxQuotaBytes: maxQuota,
		Port:          getEnv("PORT", "8080"),
		Env:           getEnv("ENV", "development"),
	}
}

// getEnv возвращает значение переменной окружения или defaultVal если она не установлена
func getEnv(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}
