// Package db handles database stuff
package db

import (
	"context"
	"database/sql"
	"log"

	"github.com/aep/parted/src/elem14"
)

// initialize tries to create the tables
func (db *Database) initialize() error {
	db.DB.Exec(SQLSchema)

	_, err := db.DB.Exec(SQLSetup)
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
	database.initialize()

	return database
}
