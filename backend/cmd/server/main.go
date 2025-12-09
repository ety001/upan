package main

import (
	"log"
	"net/http"

	"vpj/internal/config"
	"vpj/internal/handlers"
	"vpj/internal/models"
	"vpj/internal/tasks"

	"github.com/gin-gonic/gin"
)

func main() {
	// Load configuration
	if err := config.Load(); err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize database
	if err := models.InitDB(); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// Start file cleaner
	cleaner := tasks.NewCleaner()
	cleaner.Start()
	defer cleaner.Stop()

	// Setup router
	router := gin.Default()

	// CORS middleware
	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// Initialize handlers
	mainHandler := handlers.NewMainHandler()

	// Routes
	router.GET("/api/config", mainHandler.GetConfig)
	router.POST("/api/upload", mainHandler.Upload)
	router.GET("/api/file/:code", mainHandler.GetFile)

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// Start server
	addr := ":" + config.AppConfig.Port
	log.Printf("Server starting on %s", addr)
	if err := router.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

