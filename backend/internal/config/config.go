package config

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	Port          string
	DatabaseURL   string
	JWTSecret     string
	AllowedOrigin string
	DefaultLang   string
	GeoapifyKey   string
	R2AccountID   string
	R2AccessKey   string
	R2SecretKey   string
	R2BucketName  string
}

func Load() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("ℹ️  لم يتم العثور على ملف .env")
	}

	return &Config{
		Port:          getEnv("PORT", "8080"),
		DatabaseURL:   getEnv("DATABASE_URL", ""),
		JWTSecret:     getEnv("JWT_SECRET", ""),
		AllowedOrigin: getEnv("ALLOWED_ORIGINS", "http://localhost:3000,http://localhost:3001,http://localhost:3002,https://worktrack-admin.vercel.app"),
		DefaultLang:   getEnv("DEFAULT_LANG", "ar"),
		GeoapifyKey:   getEnv("GEOAPIFY_KEY", ""),
		R2AccountID:   getEnv("R2_ACCOUNT_ID", ""),
		R2AccessKey:   getEnv("R2_ACCESS_KEY", ""),
		R2SecretKey:   getEnv("R2_SECRET_KEY", ""),
		R2BucketName:  getEnv("R2_BUCKET_NAME", "worktrack-uploads"),
	}
}

// ValidateProduction prevents the API from starting with unsafe placeholder
// values when deployed. Local development keeps convenient defaults.
func (c *Config) ValidateProduction() error {
	if strings.EqualFold(getEnv("APP_ENV", "development"), "production") {
		if c.DatabaseURL == "" {
			return fmt.Errorf("DATABASE_URL must be configured in production")
		}
		if len(c.JWTSecret) < 32 || strings.Contains(strings.ToLower(c.JWTSecret), "change_this") {
			return fmt.Errorf("JWT_SECRET must be at least 32 characters and not a placeholder in production")
		}
		if c.AllowedOrigin == "" || strings.Contains(c.AllowedOrigin, "*") {
			return fmt.Errorf("ALLOWED_ORIGINS must list explicit HTTPS frontend origins in production")
		}
	}
	return nil
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
