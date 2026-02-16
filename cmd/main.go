package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jmoiron/sqlx"

	"github.com/ArkaniLoveCoding/Golang-Restfull-Api-MySql/cmd/api"
	"github.com/ArkaniLoveCoding/Golang-Restfull-Api-MySql/config"
	"github.com/ArkaniLoveCoding/Golang-Restfull-Api-MySql/db"
)

func main() {
	
	// Initialize configuration
	cfg := config.ConfigInitialize()

	// Create database configuration from config
	dbConfig := db.Config{
		Host:            cfg.PostgresHost,
		Port:            cfg.PostgresPort,
		Username:        cfg.PostgresUser,
		Password:        cfg.PostgresPassword,
		Database:        cfg.PostgresName,
		MaxOpenConns:    cfg.PostgresMaxOpenConns,
		MaxIdleConns:    cfg.PostgresMaxIdleConns,
		ConnMaxLifetime: cfg.PostgresConnMaxLifetime,
		ConnMaxIdleTime: cfg.PostgresConnMaxIdleTime,
		SSLMode:         cfg.PostgresSSLMode,
	}

	// Connect to database with retry logic
	log.Println("Connecting to PostgreSQL database...")
	database, err := db.NewConnectionWithRetry(dbConfig, cfg.DbMaxRetries, cfg.DbRetryDelay)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	log.Println(cfg.PostgresName)
	log.Println(cfg.PostgresPassword)

	defer func() {
		log.Println("Closing database connection...")
		if err := db.Close(database); err != nil {
			log.Printf("Error closing database: %v", err.Error())
		}
	}()

	// Create API server with database dependency
	server := api.ApiServerAddr(cfg.Port, database)

	// Start server in a goroutine
	go func() {
		log.Printf("Starting server on port %s...", cfg.Port)
		if err := server.Run(); err != nil {
			log.Printf("Server error: %v", err.Error())
		}
	}()

	// Wait for interrupt signal for graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// Create context with timeout for graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Shutdown server
	if err := server.Shutdown(ctx); err != nil {
		log.Printf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited properly")

}

func InitStorage(database *sqlx.DB) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := database.PingContext(ctx); err != nil {
		log.Printf("Database ping failed: %v", err)
	} else {
		log.Println("Database connection verified")
	}
}
