package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/subrat-dwi/shubserver/internal/app"
	"github.com/subrat-dwi/shubserver/internal/db"
)

func main() {
	// Load environment variables from .env file if it exists
	_ = godotenv.Load()

	// Set DATABASE_URL from environment variables if not already set
	if os.Getenv("DATABASE_URL") == "" {
		env := os.Getenv("APP_ENV")
		if env == "docker" {
			if v := os.Getenv("DATABASE_URL_DOCKER"); v != "" {
				_ = os.Setenv("DATABASE_URL", v)
			}
		} else {
			if v := os.Getenv("DATABASE_URL_LOCAL"); v != "" {
				_ = os.Setenv("DATABASE_URL", v)
			}
		}
	}

	// Get configuration from environment
	version := os.Getenv("APP_VERSION")
	if version == "" {
		version = "dev"
	}

	environment := os.Getenv("APP_ENV")
	if environment == "" {
		environment = "development"
	}

	// Connect to the database
	dbPool := db.ConnectDB()
	defer dbPool.Close()

	// Set up the application with the database connection
	s := app.Setup(dbPool, version, environment)

	fmt.Printf("Listening to http://localhost%s (version: %s, env: %s)\n", s.Addr, version, environment)
	err := http.ListenAndServe(s.Addr, s.Router)
	if err != nil && err != http.ErrServerClosed {
		fmt.Fprintf(os.Stderr, "Server error: %v\n", err)
		os.Exit(1)
	}
}
