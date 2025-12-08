package main

import (
	"biblia-am-pm/internal/database"
	"biblia-am-pm/internal/handlers"
	"biblia-am-pm/internal/middleware"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
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

	// Setup router
	r := mux.NewRouter()

	// CORS middleware - apply to all routes
	r.Use(corsMiddleware)

	// Public routes
	api := r.PathPrefix("/api").Subrouter()
	api.HandleFunc("/auth/register", authHandler.Register).Methods("POST", "OPTIONS")
	api.HandleFunc("/auth/login", authHandler.Login).Methods("POST", "OPTIONS")

	// Protected routes
	protected := api.PathPrefix("").Subrouter()
	protected.Use(corsMiddleware)
	protected.Use(middleware.AuthMiddleware)
	protected.HandleFunc("/readings/today", readingsHandler.GetTodayReadings).Methods("GET", "OPTIONS")
	protected.HandleFunc("/readings/mark-completed", readingsHandler.MarkCompleted).Methods("POST", "OPTIONS")
	protected.HandleFunc("/progress", readingsHandler.GetProgress).Methods("GET", "OPTIONS")

	// Health check
	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}).Methods("GET")

	port := os.Getenv("API_PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	if err := http.ListenAndServe(":"+port, r); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		allowedOrigins := []string{"http://localhost:3000", "http://localhost:3001"}
		
		// Allow origin if it's in the allowed list, otherwise use *
		allowOrigin := "*"
		if origin != "" {
			for _, allowed := range allowedOrigins {
				if origin == allowed {
					allowOrigin = origin
					break
				}
			}
		}
		
		w.Header().Set("Access-Control-Allow-Origin", allowOrigin)
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")
		w.Header().Set("Access-Control-Max-Age", "3600")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
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

