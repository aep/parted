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
	Manufacturer string
	PartNb       string
	Description  string
	Image        string
	Category     string
	Spec         []Attributes
	Stock        int
	Used         int
	OrderNb      string
	BarcodeID    int
}

// toItems is used to map an inbound post to a series of items
func (i InboundPOST) toItems() []Item {
	items := make([]Item, 0, len(i.Data))
	for _, item := range i.Data {
		items = append(items, Item{
			Manufacturer: item.Product.Brandname,
			PartNb:       item.Product.Translatedmanufacturerpartnumber,
			Description:  item.Product.Displayname,
			Image:        item.Product.Image.Vrntpath,
			Category:     item.Product.Displayname, // ?
			Stock:        item.Amount,
			Used:         0,
			OrderNb:      i.OrderNumber,
			BarcodeID:    item.Product.Inventorycode,
			Spec:         item.Product.Attributes,
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

	// so a schema is always available
	_, _ = db.Exec(schema)

	return &Database{DB: db}
}

// Store implements Storer
func (db *Database) Store(ctx context.Context, items []Item) error {
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

	for _, item := range items {
		exec, err := insertItemStmt.Exec(
			&item.Manufacturer,
			&item.PartNb,
			&item.Description,
			&item.Image,
			&item.Category,
			&item.Stock,
			&item.OrderNb,
			&item.BarcodeID,
		)
		if err != nil {
			return err
		}

		elementID, err := exec.LastInsertId()
		if err != nil {
			return err
		}

		for _, s := range item.Spec {
			_, err := insertMetaStmt.Exec(
				&elementID,
				&s.Attributelabel,
				&s.Attributeunit,
				&s.Attributevalue,
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
