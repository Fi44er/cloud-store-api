package config

import (
	"fmt"
	"strings"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type Config struct {
	NodeEnv string `mapstructure:"NODE_ENV"`
	Port    int    `mapstructure:"PORT"`
	AppURL  string `mapstructure:"APP_URL"`

	DBHost      string `mapstructure:"DB_HOST"`
	DBPort      int    `mapstructure:"DB_PORT"`
	DBUser      string `mapstructure:"DB_USER"`
	DBPassword  string `mapstructure:"DB_PASSWORD"`
	DBName      string `mapstructure:"DB_NAME"`
	DatabaseURL string `mapstructure:"DATABASE_URL"`

	RedisHost     string `mapstructure:"REDIS_HOST"`
	RedisPort     int    `mapstructure:"REDIS_PORT"`
	RedisPassword string `mapstructure:"REDIS_PASSWORD"`
	RedisURL      string `mapstructure:"REDIS_URL"`

	SessionSecret           string `mapstructure:"SESSION_SECRET"`
	SessionName             string `mapstructure:"SESSION_NAME"`
	SessionMaxAge           int    `mapstructure:"SESSION_MAX_AGE"`
	SessionRefreshThreshold int    `mapstructure:"SESSION_REFRESH_THRESHOLD"`
	SessionRolling          bool   `mapstructure:"SESSION_ROLLING"`

	CookieSecure   bool   `mapstructure:"COOKIE_SECURE"`
	CookieHTTPOnly bool   `mapstructure:"COOKIE_HTTP_ONLY"`
	CookieSameSite string `mapstructure:"COOKIE_SAME_SITE"`

	CORSOrigin      string `mapstructure:"CORS_ORIGIN"`
	CORSCredentials bool   `mapstructure:"CORS_CREDENTIALS"`
}

func validateConfig(config *Config) error {
	configMap := map[string]any{
		// App settings
		"NODE_ENV": config.NodeEnv,
		"PORT":     config.Port,
		"APP_URL":  config.AppURL,

		// PostgreSQL settings
		"DB_HOST":      config.DBHost,
		"DB_PORT":      config.DBPort,
		"DB_USER":      config.DBUser,
		"DB_PASSWORD":  config.DBPassword,
		"DB_NAME":      config.DBName,
		"DATABASE_URL": config.DatabaseURL,

		// Redis settings
		"REDIS_HOST": config.RedisHost,
		"REDIS_PORT": config.RedisPort,
		"REDIS_URL":  config.RedisURL,

		// Session & Auth settings
		"SESSION_SECRET":            config.SessionSecret,
		"SESSION_NAME":              config.SessionName,
		"SESSION_MAX_AGE":           config.SessionMaxAge,
		"SESSION_REFRESH_THRESHOLD": config.SessionRefreshThreshold,
		"SESSION_ROLLING":           config.SessionRolling,

		// Cookie security
		"COOKIE_SECURE":    config.CookieSecure,
		"COOKIE_HTTP_ONLY": config.CookieHTTPOnly,
		"COOKIE_SAME_SITE": config.CookieSameSite,

		// CORS settings
		"CORS_ORIGIN":      config.CORSOrigin,
		"CORS_CREDENTIALS": config.CORSCredentials,
	}

	for key, value := range configMap {
		if isEmptyValue(value) {
			return fmt.Errorf("missing required configuration field: %s", key)
		}
	}

	return nil
}

func LoadConfig(path string) (config *Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName(".env")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return nil, err
	}

	err = viper.Unmarshal(&config)
	if err != nil {
		return nil, err
	}

	err = validateConfig(config)
	if err != nil {
		return nil, err
	}

	return config, nil
}

func Load(path string) error {
	err := godotenv.Load(path)
	if err != nil {
		return err
	}

	return nil
}

func isEmptyValue(value any) bool {
	switch v := value.(type) {
	case string:
		return strings.TrimSpace(v) == ""
	case int:
		return v == 0
	case int64:
		return v == 0
	case bool:
		return false // bool значения не проверяем на пустоту, так как они всегда имеют значение
	case nil:
		return true
	default:
		return false
	}
}
