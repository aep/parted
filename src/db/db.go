// Package db handles database stuff
package db

import (
	"context"
	"database/sql"
	"log"

	"github.com/aep/parted/src/elem14"
)

// initializeSchema tries to create the tables
func (db *Database) initializeSchema() error {
	_, err := db.DB.Exec(SQLSchema)
	if err != nil {
		return err
	}
	return err
}

// Storer stores items inside the database
type Storer interface {
	Store(context.Context, []elem14.Item) error
}

// Database holds anything relevent to the database
type Database struct {
	DB *sql.DB
}

// Connect to the databse and initialize the schema
func Connect() *Database {
	db, err := sql.Open("sqlite3", "./items.db")
	if err != nil {
		log.Fatal(err)
	}

	database := &Database{DB: db}

	// so a schema is always available
	//nolint
	database.initializeSchema()

	return database
}
