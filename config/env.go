package config

import (
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type ConfigParams struct {
	PublicHost      string
	Port 			string
	// PostgreSQL specific
	PostgresHost            string
	PostgresPort            string
	PostgresName            string
	PostgresUser            string
	PostgresPassword        string
	PostgresMaxOpenConns    int
	PostgresMaxIdleConns    int
	PostgresConnMaxLifetime time.Duration
	PostgresConnMaxIdleTime time.Duration
	PostgresSSLMode         string
	// Retry settings
	DbMaxRetries  int
	DbRetryDelay  time.Duration
}

func ConfigInitialize() ConfigParams {

	_ = godotenv.Load()

	return ConfigParams{
		PublicHost: KeyEnvLookUp("PUBLIC_HOST", "http://localhost"),
		Port	  : KeyEnvLookUp("PORT", ":8080"),
		// PostgreSQL specific config
		PostgresHost:            KeyEnvLookUp("DB_HOST", "localhost"),
		PostgresPort:            KeyEnvLookUp("DB_PORT", "5432"),
		PostgresName:            KeyEnvLookUp("DB_NAME", "School-manajement"),
		PostgresUser:            KeyEnvLookUp("DB_USER", "appuser2"),
		PostgresPassword:        KeyEnvLookUp("DB_PASSWORD", "app123"),
		PostgresMaxOpenConns:    getEnvInt("DB_MAX_OPEN_CONNS", 25),
		PostgresMaxIdleConns:    getEnvInt("DB_MAX_IDLE_CONNS", 5),
		PostgresConnMaxLifetime: getEnvDuration("DB_CONN_MAX_LIFETIME", 5*time.Minute),
		PostgresConnMaxIdleTime: getEnvDuration("DB_CONN_MAX_IDLE_TIME", 5*time.Minute),
		PostgresSSLMode:         KeyEnvLookUp("DB_SSL_MODE", "disable"),
		// Retry settings
		DbMaxRetries: getEnvInt("DB_MAX_RETRIES", 3),
		DbRetryDelay: getEnvDuration("DB_RETRY_DELAY", 5*time.Second),
	}

}

func KeyEnvLookUp(key string, fallback string) string {

	if value, ok := os.LookupEnv(key); ok {
		return value
	}

	return fallback
}

func getEnvInt(key string, fallback int) int {
	if value, ok := os.LookupEnv(key); ok {
		if intVal, err := strconv.Atoi(value); err == nil {
			return intVal
		}
	}
	return fallback
}

func getEnvDuration(key string, fallback time.Duration) time.Duration {
	if value, ok := os.LookupEnv(key); ok {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return fallback
}
