package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"

	"github.com/ArkaniLoveCoding/Shcool-manajement/cmd/api"
	"github.com/ArkaniLoveCoding/Shcool-manajement/config"
	"github.com/ArkaniLoveCoding/Shcool-manajement/db"
	"github.com/ArkaniLoveCoding/Shcool-manajement/middleware/logger"
)

func main() {

	// Initialize logger FIRST (before anything else)
	logger.Init()
	defer logger.Sync()

	logger.Log.Info("Starting application...")

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
	logger.Log.Info("Connecting to PostgreSQL database...")
	database, err := db.NewConnectionWithRetry(dbConfig, cfg.DbMaxRetries, cfg.DbRetryDelay)
	if err != nil {
		logger.Log.Fatal("Failed to connect to database",
			zap.Error(err),
		)
		log.Fatalf("Failed to connect to database: %v", err)
	}

	logger.Log.Info("Database connected successfully",
		zap.String("database", cfg.PostgresName),
	)

	defer func() {
		logger.Log.Info("Closing database connection...")
		if err := db.Close(database); err != nil {
			logger.Log.Error("Error closing database",
				zap.Error(err),
			)
			log.Printf("Error closing database: %v", err.Error())
		}
	}()

	// Create API server with database dependency
	// Logger is already initialized at this point
	server := api.ApiServerAddr(cfg.Port, database)

	// Start server in a goroutine
	go func() {
		logger.Log.Info("Starting server",
			zap.String("port", cfg.Port),
		)
		if err := server.Run(); err != nil {
			logger.Log.Error("Server error",
				zap.Error(err),
			)
			log.Printf("Server error: %v", err.Error())
		}
	}()

	// Wait for interrupt signal for graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Log.Info("Shutting down server...")

	// Create context with timeout for graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Shutdown server
	if err := server.Shutdown(ctx); err != nil {
		logger.Log.Error("Server forced to shutdown",
			zap.Error(err),
		)
		log.Printf("Server forced to shutdown: %v", err)
	}

	logger.Log.Info("Server exited properly")
	log.Println("Server exited properly")

}

func InitStorage(database *sqlx.DB) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := database.PingContext(ctx); err != nil {
		logger.Log.Error("Database ping failed",
			zap.Error(err),
		)
		log.Printf("Database ping failed: %v", err)
	} else {
		logger.Log.Info("Database connection verified")
	}
}

