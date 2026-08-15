package app

import "os"

type Config struct {
	Addr          string
	DSN           string
	SessionSecret string
	UploadDir     string
}

func LoadConfig() Config {
	return Config{
		Addr:          env("RBAC_ADDR", ":8080"),
		DSN:           env("RBAC_DSN", "root:123321qQ@tcp(127.0.0.1:3306)/shop?charset=utf8mb4&parseTime=True&loc=Local"),
		SessionSecret: env("RBAC_SESSION_SECRET", "change-me-before-production"),
		UploadDir:     env("RBAC_UPLOAD_DIR", "web/static/uploads"),
	}
}

func env(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
