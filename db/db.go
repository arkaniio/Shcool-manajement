package db

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

// Config holds database connection configuration
type Config struct {
	Host            string
	Port            string
	Username        string
	Password        string
	Database        string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
	ConnMaxIdleTime time.Duration
	SSLMode         string
}

// DefaultConfig returns default PostgreSQL configuration
func DefaultConfig() Config {
	return Config{
		Host:            "localhost",
		Port:            "5432",
		Username:        "appuser2",
		Password:        "app123",
		Database:        "School-manajement",
		MaxOpenConns:    25,
		MaxIdleConns:    5,
		ConnMaxLifetime: 5 * time.Minute,
		ConnMaxIdleTime: 5 * time.Minute,
		SSLMode:         "disable",
	}
}

// NewConnection creates a new PostgreSQL database connection with connection pool
func NewConnection(cfg Config) (*sqlx.DB, error) {
	// Build DSN (Data Source Name)
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host,
		cfg.Port,
		cfg.Username,
		cfg.Password,
		cfg.Database,
		cfg.SSLMode,
	)

	// Connect to database
	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Configure connection pool
	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetMaxIdleConns(cfg.MaxIdleConns)
	db.SetConnMaxLifetime(cfg.ConnMaxLifetime)
	db.SetConnMaxIdleTime(cfg.ConnMaxIdleTime)

	// Verify connection with ping
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	log.Println("Successfully connected to PostgreSQL database")
	return db, nil
}

// NewConnectionWithRetry creates a connection with automatic retry logic
func NewConnectionWithRetry(cfg Config, maxRetries int, retryDelay time.Duration) (*sqlx.DB, error) {
	var db *sqlx.DB
	var err error

	for i := 0; i < maxRetries; i++ {
		db, err = NewConnection(cfg)
		if err == nil {
			return db, nil
		}

		log.Printf("Failed to connect to database (attempt %d/%d): %v", i+1, maxRetries, err)
		
		if i < maxRetries-1 {
			time.Sleep(retryDelay * time.Duration(i+1)) // Exponential backoff
		}
	}

	return nil, fmt.Errorf("failed to connect after %d retries: %w", maxRetries, err)
}

// HealthCheck verifies database connectivity
func HealthCheck(ctx context.Context, db *sqlx.DB) error {
	return db.PingContext(ctx)
}

// Close gracefully closes the database connection
func Close(db *sqlx.DB) error {
	if db != nil {
		return db.Close()
	}
	return nil
}

