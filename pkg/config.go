package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	DBHost        string `mapstructure:"DB_HOST"`
	DBPort        int    `mapstructure:"DB_PORT"`
	DBUser        string `mapstructure:"DB_USER"`
	DBPassword    string `mapstructure:"DB_PASSWORD"`
	DBName        string `mapstructure:"DB_NAME"`
	DBSSLMode     string `mapstructure:"DB_SSL_MODE"`
	DBAutoMigrate bool   `mapstructure:"DB_AUTO_MIGRATE"`

	RedisAddress  string `mapstructure:"REDIS_ADDRESS"`
	RedisPassword string `mapstructure:"REDIS_PASSWORD"`
	RedisDB       int    `mapstructure:"REDIS_DB"`

	WeatherAPIKey string `mapstructure:"WEATHER_API_KEY"`

	SMTPHost     string `mapstructure:"SMTP_HOST"`
	SMTPPort     int    `mapstructure:"SMTP_PORT"`
	SMTPUsername string `mapstructure:"SMTP_USERNAME"`
	SMTPPassword string `mapstructure:"SMTP_PASSWORD"`
	SMTPFrom     string `mapstructure:"SMTP_FROM"`

	BaseURL    string `mapstructure:"BASE_URL"`
	SwaggerURL string `mapstructure:"SWAGGER_URL"`

	APP_Host string `mapstructure:"APP_HOST"`
	APP_Port int    `mapstructure:"APP_PORT"`
}

func (c *Config) GetDSN() string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		c.DBHost, c.DBPort, c.DBUser, c.DBPassword, c.DBName, c.DBSSLMode,
	)
}

func LoadConfig(configPath string) (*Config, error) {
	v := viper.New()

	v.AutomaticEnv()

	v.SetDefault("DB_HOST", "localhost")
	v.SetDefault("DB_PORT", 5432)
	v.SetDefault("DB_USER", "postgres")
	v.SetDefault("DB_NAME", "weather_app")
	v.SetDefault("DB_SSL_MODE", "disable")

	v.SetDefault("REDIS_ADDRESS", "localhost:6379")
	v.SetDefault("REDIS_DB", 0)

	v.SetDefault("BASE_URL", "http://localhost:8080")
	v.SetDefault("SWAGGER_URL", "http://localhost:8080/swagger/doc.json")

	v.SetDefault("APP_HOST", "localhost")
	v.SetDefault("APP_PORT", 8080)

	if configPath != "" {
		v.AddConfigPath(configPath)
		v.SetConfigName(".env")
		v.SetConfigType("env")

		if err := v.ReadInConfig(); err != nil {
			if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
				return nil, fmt.Errorf("error reading config file: %w", err)
			}
		}
	}

	var config Config

	if err := v.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("error unmarshaling config: %w", err)
	}

	if err := validateConfig(&config); err != nil {
		return nil, err
	}

	return &config, nil
}

func validateConfig(config *Config) error {
	missingFields := []string{}

	if config.WeatherAPIKey == "" {
		missingFields = append(missingFields, "WEATHER_API_KEY")
	}

	if len(missingFields) > 0 {
		return fmt.Errorf("missing required configuration fields: %s", strings.Join(missingFields, ", "))
	}

	return nil
}
