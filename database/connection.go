package database

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3" // SQLite driver
)

// InitDB initializes the SQLite database connection
func InitDB() *sql.DB {
	db, err := sql.Open("sqlite3", "./test.db")
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Create table with additional fields if it doesn't exist
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		email TEXT NOT NULL UNIQUE,
		password TEXT NOT NULL,
		phone TEXT,
		address TEXT,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	)`)
	if err != nil {
		log.Fatalf("Failed to create table: %v", err)
	}

	return db
}
