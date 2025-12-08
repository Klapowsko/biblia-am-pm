package main

import (
	"biblia-am-pm/internal/database"
	"biblia-am-pm/internal/handlers"
	"biblia-am-pm/internal/middleware"
	"log"
	"os"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize database
	if err := database.InitDB(); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer database.CloseDB()

	// Run migrations
	if err := runMigrations(); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	// Initialize handlers
	authHandler := handlers.NewAuthHandler()
	readingsHandler := handlers.NewReadingsHandler()

	// Setup Gin router
	r := gin.Default()

	// CORS middleware
	config := cors.DefaultConfig()
	// Get allowed origins from environment or use defaults
	allowedOrigins := os.Getenv("CORS_ALLOWED_ORIGINS")
	if allowedOrigins == "" {
		// Default origins: localhost and common server domains
		config.AllowOrigins = []string{
			"http://localhost:3001",
			"http://hiagoserver.local:3001",
			"http://hiagoserver.local",
		}
	} else {
		// Parse comma-separated origins from environment
		origins := []string{}
		for _, origin := range strings.Split(allowedOrigins, ",") {
			origins = append(origins, strings.TrimSpace(origin))
		}
		config.AllowOrigins = origins
	}
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"Content-Type", "Authorization", "X-Requested-With"}
	config.AllowCredentials = true
	r.Use(cors.New(config))

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.String(200, "OK")
	})

	// API routes
	api := r.Group("/api")
	{
		// Public routes
		api.POST("/auth/register", authHandler.Register)
		api.POST("/auth/login", authHandler.Login)
	}

	// Protected routes - create separate group with auth middleware
	protected := r.Group("/api")
	protected.Use(middleware.AuthMiddleware())
	{
		protected.GET("/readings/today", readingsHandler.GetTodayReadings)
		protected.POST("/readings/mark-completed", readingsHandler.MarkCompleted)
		protected.GET("/progress", readingsHandler.GetProgress)
	}

	port := os.Getenv("API_PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}


func runMigrations() error {
	migrationSQL := `
	-- Create users table
	CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		email VARCHAR(255) UNIQUE NOT NULL,
		password VARCHAR(255) NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);

	-- Create reading_plans table
	CREATE TABLE IF NOT EXISTS reading_plans (
		id SERIAL PRIMARY KEY,
		day_of_year INTEGER NOT NULL UNIQUE,
		old_testament_ref VARCHAR(255),
		new_testament_ref VARCHAR(255),
		psalms_ref VARCHAR(255),
		proverbs_ref VARCHAR(255)
	);

	-- Create user_progress table
	CREATE TABLE IF NOT EXISTS user_progress (
		id SERIAL PRIMARY KEY,
		user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
		reading_plan_id INTEGER NOT NULL REFERENCES reading_plans(id) ON DELETE CASCADE,
		date DATE NOT NULL,
		morning_completed BOOLEAN DEFAULT FALSE,
		evening_completed BOOLEAN DEFAULT FALSE,
		completed_at TIMESTAMP,
		UNIQUE(user_id, reading_plan_id, date)
	);

	-- Create indexes for better performance
	CREATE INDEX IF NOT EXISTS idx_user_progress_user_id ON user_progress(user_id);
	CREATE INDEX IF NOT EXISTS idx_user_progress_date ON user_progress(date);
	CREATE INDEX IF NOT EXISTS idx_reading_plans_day_of_year ON reading_plans(day_of_year);
	`

	_, err := database.DB.Exec(migrationSQL)
	return err
}

