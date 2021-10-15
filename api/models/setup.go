// Package models contains the models used by the API
package models

import (
	"fmt"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// OpenDatabase opens an SQLite database file on the given path. It auto-migrates the models.
func OpenDatabase(path string) (*gorm.DB, error) {
	database, err := gorm.Open(sqlite.Open(path), &gorm.Config{})

	if err != nil {
		return nil, fmt.Errorf("could not open database: %w", err)
	}

	database.Exec("PRAGMA synchronous = OFF")
	database.Exec("PRAGMA journal_mode = WAL")
	database.Exec("PRAGMA foreign_keys = ON")

	database.AutoMigrate(&User{})
	database.AutoMigrate(&HighScore{})

	return database, nil
}
