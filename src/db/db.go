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
	ID           int         `db:"id" json:"id"`
	Manufacturer string      `db:"manufacturer" json:"manufacturer"`
	PartNumber   string      `db:"part_number" json:"part_number"`
	Description  string      `db:"description" json:"description"`
	Image        string      `db:"image" json:"image"`
	Stock        int         `db:"stock" json:"stock"`
	Used         int         `db:"used" json:"used"`
	OrderNumber  string      `db:"order_number" json:"order_number"`
	BarcodeID    int         `db:"barcode_id" json:"barcode_id"`
	InsertDate   time.Time   `db:"insert_date" json:"insert_date"`
	Attributes   []Attribute `json:"attributes"`
}

type Attribute struct {
	Label string `json:"label"`
	Value string `json:"value"`
	Unit  string `json:"unit"`
}

// IntermediateAttr holds the intermediate values to handle the conversions
type IntermediateAttr struct {
	Label string         `db:"label"`
	Value sql.NullString `db:"value"`
	Unit  sql.NullString `db:"unit"`
}
