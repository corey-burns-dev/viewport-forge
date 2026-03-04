package config

import "os"

type Config struct {
	APIPort       string
	RedisAddr     string
	QueueKey      string
	StatusPrefix  string
	AllowedOrigin string
}

func Load() Config {
	return Config{
		APIPort:       getEnv("API_PORT", "8080"),
		RedisAddr:     getEnv("REDIS_ADDR", "localhost:6379"),
		QueueKey:      getEnv("QUEUE_KEY", "vf:capture_jobs"),
		StatusPrefix:  getEnv("STATUS_PREFIX", "vf:capture_status:"),
		AllowedOrigin: getEnv("ALLOWED_ORIGIN", "http://localhost:5173"),
	}
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
