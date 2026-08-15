package app

import "os"

type Config struct {
	Addr      string
	DSN       string
	UploadDir string
}

func LoadConfig() Config {
	return Config{
		Addr:      env("FRONT_API_ADDR", ":8090"),
		DSN:       env("FRONT_DSN", "root:123321qQ@tcp(127.0.0.1:3306)/shop?charset=utf8mb4&parseTime=True&loc=Local"),
		UploadDir: env("FRONT_UPLOAD_DIR", "../admin/web/static/uploads"),
	}
}

func env(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
