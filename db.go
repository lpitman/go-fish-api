package main

import (
	"database/sql"
	"log"
	"os"
	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func ConnectDB() error {
	dbPath := os.Getenv("SQLITE_DB_PATH")
	if dbPath == "" {
		dbPath = "fish_tracking.db"
	}

	var err error
	DB, err = sql.Open("sqlite3", dbPath)
	if err != nil {
		return err
	}
	
	if err = DB.Ping(); err != nil {
		return err;
	}

	log.Println("Connected to SQLite database successfully at:", dbPath	)

	// Initialize our fishy table if it doesn't exist
	createTableQuery := `
	CREATE TABLE IF NOT EXISTS fish (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		species TEXT NOT NULL,
		tracking_info TEXT NOT NULL,
		weight_kg REAL NOT NULL
	);`
	_, err = DB.Exec(createTableQuery)
	if err != nil {
		return err
	}

	log.Println("Fish table is ready.")

	return nil
}
