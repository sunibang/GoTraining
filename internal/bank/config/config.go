package config

import (
	"os"
)

type Config struct {
	DatabaseURL  string
	JWTSecret    string
	ServiceName  string
	OTelEndpoint string
	Port         string
}

var Values Config

func Init() {
	Values = Config{
		DatabaseURL:  envOrDefault("DATABASE_URL", "postgres://gotrainer:verysecret@localhost:5432/gobank?sslmode=disable"),
		JWTSecret:    envOrDefault("JWT_SECRET", "super-secret-for-training-only"),
		ServiceName:  envOrDefault("SERVICE_NAME", "bank-api"),
		OTelEndpoint: envOrDefault("OTEL_EXPORTER_OTLP_ENDPOINT", "localhost:4318"),
		Port:         envOrDefault("PORT", "8080"),
	}
}

func envOrDefault(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
