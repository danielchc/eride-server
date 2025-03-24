package db

import (
	"database/sql"
	"fmt"
	"os"

	_ "modernc.org/sqlite"
)

type SQLiteDB struct {
	DB *sql.DB
}

func NewSQLiteDB(dbPath string) (*SQLiteDB, error) {
	// Ensure the database file exists or create it
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		// If the database doesn't exist, SQLite will create it automatically when opened
		_, err := os.Create(dbPath) // Create the database file if it does not exist
		if err != nil {
			return nil, fmt.Errorf("failed to create database file: %v", err)
		}
	}

	// Open the SQLite database (it will be created if it doesn't exist)
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	// Create the users table if it doesn't already exist
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			username TEXT PRIMARY KEY,
			password TEXT,
			data_path TEXT
		)
	`)

	if err != nil {
		return nil, fmt.Errorf("failed to create users table: %v", err)
	}

	return &SQLiteDB{
		DB: db,
	}, nil
}
