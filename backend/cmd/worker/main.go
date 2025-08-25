package main

import (
	"log"
	"os"
	"os/signal"
	"pickup-queue/internal/repository"
	"pickup-queue/internal/usecase"
	"pickup-queue/pkg/database"
	"pickup-queue/pkg/logger"
	"syscall"
	"time"

	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	// Initialize logger
	appLogger := logger.New()

	// Initialize database
	dbConfig := database.GetConfigFromEnv()
	db, err := database.NewConnection(dbConfig)
	if err != nil {
		appLogger.Error("Failed to connect to database:", err)
		os.Exit(1)
	}

	// Initialize repositories
	packageRepo := repository.NewPackageRepository(db)

	// Initialize use cases
	packageUsecase := usecase.NewPackageUsecase(packageRepo)

	// Create a ticker that runs every hour
	ticker := time.NewTicker(1 * time.Hour)
	defer ticker.Stop()

	// Channel to listen for interrupt signal
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	appLogger.Info("Package expiry worker started")

	// Run initial check
	if err := packageUsecase.MarkExpiredPackages(); err != nil {
		appLogger.Error("Error marking expired packages:", err)
	} else {
		appLogger.Info("Initial expired packages check completed")
	}

	for {
		select {
		case <-ticker.C:
			appLogger.Info("Running expired packages check...")
			if err := packageUsecase.MarkExpiredPackages(); err != nil {
				appLogger.Error("Error marking expired packages:", err)
			} else {
				appLogger.Info("Expired packages check completed")
			}
		case <-interrupt:
			appLogger.Info("Shutting down worker...")
			return
		}
	}
}
