package config // File belongs to config package

import (
	"fmt"     // String formatting
	"log"     // Print logs/messages
	"os"      // Access operating system and env variables
	"strconv" // Convert string ↔ number

	"github.com/joho/godotenv" // Read values from .env file and set them as environment variables.
)

// Config holds every setting our application needs.
// Think of it as a settings panel — everything in one place.

type Config struct {
	AppPort        string
	DBHost         string
	DBPort         string
	DBUser         string
	DBPassword     string
	DBName         string
	DBURL          string
	JWTSecret      string
	JWTExpiryHours int
}

// LoadConfig reads from .env file(for local dev) and environment variables, then fills the Config struct.
// In Production(K8s), values come from ConfigMaps and Secrets, which are set as environment variables - not .env file.

func LoadConfig() *Config {
	// Function name = LoadConfig
	// Return type = *Config

	// In local development, try to load from .env file.
	// In Docker/production, environment variables are already passed in, so .env is optional.
	if _, err := os.Stat(".env"); err == nil {
		if err := godotenv.Load(); err != nil {
			log.Println("Failed to load .env file, using environment variables instead.")
		}
	}

	jwtExpiry, err := strconv.Atoi(getEnv("JWT_EXPIRY_HOURS", "24")) // Convert JWT_EXPIRY_HOURS from string to int

	// os.Getenv() accepts only one argument
	if err != nil {
		jwtExpiry = 24 // Default to 24 hrs if conversion fails

	}

	cfg := &Config{
		AppPort:        getEnv("APP_PORT", "8080"),
		DBHost:         getEnv("DB_HOST", "localhost"),
		DBPort:         getEnv("DB_PORT", "5432"),
		DBUser:         getEnv("DB_USER", "postgres"),
		DBPassword:     getEnv("DB_PASSWORD", "postgres"),
		DBName:         getEnv("DB_NAME", "taskmanager"),
		JWTSecret:      getEnv("JWT_SECRET", "your-super-secret-key-change-in-production"),
		JWTExpiryHours: jwtExpiry,
	}

	// Build the PostgreSQL connection string
	cfg.DBURL = fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName,
	)
	return cfg
}

// getEnv is a helper function to read an environment variable or return a default value if not set.
func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
