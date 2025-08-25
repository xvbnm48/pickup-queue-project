package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/lib/pq"
)

type Config struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

func NewConnection(config *Config) (*sql.DB, error) {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		config.Host, config.Port, config.User, config.Password, config.DBName, config.SSLMode)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Test the connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return db, nil
}

func GetConfigFromEnv() *Config {
	return &Config{
		Host:     getEnv("DB_HOST", "172.18.125.255"),
		Port:     getEnv("DB_PORT", "5432"),
		User:     getEnv("DB_USER", "postgres"),
		Password: getEnv("DB_PASSWORD", "vini"),
		DBName:   getEnv("DB_NAME", "pickup_queue"),
		SSLMode:  getEnv("DB_SSL_MODE", "disable"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// LogQuery logs SQL queries with execution time
func LogQuery(query string, args []interface{}, startTime time.Time) {
	duration := time.Since(startTime)
	log.Printf("[SQL Query] Duration: %v | Query: %s | Args: %v", duration, query, args)
}

// LogQueryError logs SQL query errors
func LogQueryError(query string, args []interface{}, err error, startTime time.Time) {
	duration := time.Since(startTime)
	log.Printf("[SQL Error] Duration: %v | Query: %s | Args: %v | Error: %v", duration, query, args, err)
}
