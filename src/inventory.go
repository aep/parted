package src

import (
	"context"
	"database/sql"
	"log"

	// database driver initialization
	_ "github.com/mattn/go-sqlite3"
)

// Item represents is an item in the inventory
type Item struct {
	Manufacturer string       `json:"manufacturer"`
	PartNumber   string       `json:"part_number"`
	Description  string       `json:"description"`
	Image        string       `json:"image"`
	Category     string       `json:"category"`
	Attributes   []Attributes `json:"attributes"`
	Stock        int          `json:"stock"`
	Used         int          `json:"used"`
	OrderNb      string       `json:"order_number"`
	BarcodeID    int          `json:"barcode_id"`
}

// toItems is used to map an inbound post to a series of items
func (e *ManufacturerPartNumberSearch) toItems() []Item {
	items := make([]Item, 0, len(e.Manufacturerpartnumbersearchreturn.Products))
	for _, item := range e.Manufacturerpartnumbersearchreturn.Products {
		items = append(items, Item{
			Manufacturer: item.Brandname,
			PartNumber:   item.Translatedmanufacturerpartnumber,
			Description:  item.Displayname,
			Image:        item.Image.Vrntpath + item.Image.Basename,
			Category:     item.Displayname, // ?
			Stock:        item.Inv,
			BarcodeID:    item.Inventorycode,
			Attributes:   item.Attributes,
		})
	}

	return items
}

// Storer stores items inside the database
type Storer interface {
	Store(context.Context, []Item) error
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

// Store implements Storer
func (db *Database) Store(ctx context.Context, items InboundPOST) error {
	tx, err := db.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	//nolint
	defer tx.Rollback()

	insertItemStmt, err := tx.Prepare(insertItem)
	if err != nil {
		return err
	}

	insertMetaStmt, err := tx.Prepare(insertAttr)
	if err != nil {
		return err
	}

	for _, item := range items.Data {
		exec, err := insertItemStmt.Exec(
			&item.Manufacturer,
			&item.PartNumber,
			&item.Description,
			&item.Image,
			&item.Category,
			&item.Stock,
			&items.OrderNumber,
			&item.BarcodeID,
		)
		if err != nil {
			return err
		}

		elementID, err := exec.LastInsertId()
		if err != nil {
			return err
		}

		for _, spec := range item.Attributes {
			_, err := insertMetaStmt.Exec(
				&elementID,
				&spec.Attributelabel,
				&spec.Attributevalue,
				&spec.Attributeunit,
			)
			if err != nil {
				return err
			}
		}

	}

	return tx.Commit()
}

// initializeSchema tries to create the tables
func (db *Database) initializeSchema() error {
	_, err := db.DB.Exec(schema)
	if err != nil {
		return err
	}
	return err
}

const schema = `
CREATE TABLE inventory (
    id INTEGER NOT NULL PRIMARY KEY,
    manufacturer VARCHAR,
    part_number VARCHAR,
    description VARCHAR,
    image VARCHAR,
    category VARCHAR,
    stock INTEGER,
    used INTEGER,
    order_number VARCHAR,
    barcode_id INTEGER
);

CREATE TABLE specifications (
    id_element INTEGER,
    label VARCHAR,
    value VARCHAR,
    unit VARCHAR,
    FOREIGN KEY(id_element) REFERENCES inventory(id)
);
`

const insertItem = `
INSERT INTO inventory (
	manufacturer,
	part_number,
	description,
	image,
	category,
	stock,
	order_number,
	barcode_id
)
VALUES (?, ?, ?, ?, ?, ?, ?, ?);
`

const insertAttr = `
INSERT INTO specifications (id_element, label, value, unit)
VALUES (?, ?, ?, ?);
`
