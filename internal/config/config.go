package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Env struct {
	AppPort     string
	DBHost      string
	DBPort      string
	DBUser      string
	DBPass      string
	DBName      string
	DBSSLMode   string
	JWTSecret   string
	AppBaseURL  string
	RabbitURL   string
	MailpitHost string
	MailpitPort string
}

func Load() Env {
	// loads .env from project root
	_ = godotenv.Load()

	env := Env{
		AppPort:    get("APP_PORT", "8080"),
		DBHost:     get("DB_HOST", "localhost"),
		DBPort:     get("DB_PORT", "5432"),
		DBUser:     get("DB_USER", "postgres"),
		DBPass:     os.Getenv("DB_PASSWORD"), // no default for secrets
		DBName:     get("DB_NAME", "qstack"),
		DBSSLMode:  get("DB_SSLMODE", "disable"),
		JWTSecret:  os.Getenv("JWT_SECRET"),
		AppBaseURL: get("APP_BASE_URL", "http://localhost:8080"),
		RabbitURL:   get("RABBITMQ_URL", "amqp://guest:guest@localhost:5672/"),
		MailpitHost: get("MAILPIT_HOST", "localhost"),
		MailpitPort: get("MAILPIT_PORT", "1025"),
	}

	if env.DBPass == "" {
		log.Println("DB_PASSWORD is empty (check your .env)")
	}
	return env
}

func get(key, fallback string) string {
	v := os.Getenv(key)
	if v == "" {
		return fallback
	}
	return v
}