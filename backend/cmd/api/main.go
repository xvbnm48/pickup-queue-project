package main

import (
	"log"
	"os"
	"pickup-queue/internal/handler"
	"pickup-queue/internal/middleware"
	"pickup-queue/internal/repository"
	"pickup-queue/internal/usecase"
	"pickup-queue/pkg/database"
	"pickup-queue/pkg/logger"

	"github.com/gin-gonic/gin"
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

	// Auto migrate
	if err := database.AutoMigrate(db); err != nil {
		appLogger.Error("Failed to migrate database:", err)
		os.Exit(1)
	}

	// Initialize repositories
	packageRepo := repository.NewPackageRepository(db)

	// Initialize use cases
	packageUsecase := usecase.NewPackageUsecase(packageRepo)

	// Initialize handlers
	packageHandler := handler.NewPackageHandler(packageUsecase)

	// Initialize Gin router
	router := gin.New()

	// Add middleware
	router.Use(middleware.Logger())
	router.Use(middleware.CORS())
	router.Use(middleware.RequestID())
	router.Use(gin.Recovery())

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "healthy",
			"service": "pickup-queue-api",
		})
	})

	// API routes
	v1 := router.Group("/api/v1")
	{
		packages := v1.Group("/packages")
		{
			packages.POST("", packageHandler.CreatePackage)
			packages.GET("", packageHandler.ListPackages)
			packages.GET("/stats", packageHandler.GetPackageStats)
			packages.GET("/:id", packageHandler.GetPackage)
			packages.GET("/order/:orderRef", packageHandler.GetPackageByOrderRef)
			packages.PATCH("/:id/status", packageHandler.UpdatePackageStatus)
			packages.DELETE("/:id", packageHandler.DeletePackage)
		}
	}

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	appLogger.Info("Starting server on port", port)
	if err := router.Run(":" + port); err != nil {
		appLogger.Error("Failed to start server:", err)
		os.Exit(1)
	}
}
