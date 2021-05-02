// Package db handles database stuff
package db

import (
	"context"
	"database/sql"
	"log"
	"time"

	"github.com/aep/parted/src/elem14"
)

// initialize tries to create the tables
func (db *Database) initialize() error {
	_, err := db.DB.Exec(SQLSchema)
	if err != nil {
		return err
	}

	_, err = db.DB.Exec(SQLSetup)
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

type Item struct {
	ID           int       `db:"id"`
	Manufacturer string    `db:"manufacturer"`
	PartNumber   string    `db:"part_number"`
	Description  string    `db:"description"`
	Image        string    `db:"image"`
	Stock        int       `db:"stock"`
	Used         int       `db:"used"`
	OrderNumber  string    `db:"order_number"`
	BarcodeID    int       `db:"barcode_id"`
	InsertDate   time.Time `db:"insert_date"`
	Attributes   []Attribute
}

type Attribute struct {
	Label string
	Value string
	Unit  string
}

// IntermediateAttr holds the intermediate values to handle the conversions
type IntermediateAttr struct {
	Label string         `db:"label"`
	Value sql.NullString `db:"value"`
	Unit  sql.NullString `db:"unit"`
}
