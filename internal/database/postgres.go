package database

import (
	"log"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

// NewPostgresDB creates and returns a database connection pool.
func NewPostgresDB(dsn string) *sqlx.DB {
	db, err := sqlx.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("Failed to open DB connection: %v", err)
	}

	// Connection pool settings — tune these based on your load
	db.SetMaxOpenConns(25)         
	db.SetMaxIdleConns(10)          
	db.SetConnMaxLifetime(5 * time.Minute) 

	// Ping verifies the connection is actually working
	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to connect to DB: %v", err)
	}

	log.Println("Connected to PostgreSQL")
	return db
}