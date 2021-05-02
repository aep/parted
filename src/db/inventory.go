package db

import (

	// database driver initialization
	"fmt"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

// StoreInbound implements Storer
func (db *Database) StoreInbound(items []Item, order string) error {
	tx, err := db.DB.Begin()
	if err != nil {
		return err
	}

	//nolint
	defer tx.Rollback()

	insertItemStmt, err := tx.Prepare(SQLInsertItem)
	if err != nil {
		return fmt.Errorf("error preparing insert item %w", err)
	}

	insertMetaStmt, err := tx.Prepare(SQLInsertAttr)
	if err != nil {
		return fmt.Errorf("error preparing insert attr %w", err)
	}

	for _, item := range items {
		exec, err := insertItemStmt.Exec(
			&item.Manufacturer,
			&item.PartNumber,
			&item.Description,
			&item.Image,
			&item.Stock,
			&order,
			&item.BarcodeID,
			time.Now().Unix(),
		)
		if err != nil {
			return fmt.Errorf("error inserting item %w", err)
		}

		elementID, err := exec.LastInsertId()
		if err != nil {
			return fmt.Errorf("error getting last id %w", err)
		}

		for _, spec := range item.Attributes {
			_, err := insertMetaStmt.Exec(
				&elementID,
				&spec.Label,
				&spec.Value,
				&spec.Unit,
			)
			if err != nil {
				return fmt.Errorf("error inserting attr %w", err)
			}
		}

	}

	return tx.Commit()
}

// ReadAll the items
func (db *Database) ReadAll() ([]Item, error) {
	rows, err := db.DB.Query(SQLReadItems)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Item
	for rows.Next() {
		var item Item
		var attr IntermediateAttr
		if err := rows.Scan(
			&item.ID,
			&item.Manufacturer,
			&item.PartNumber,
			&item.Description,
			&item.Image,
			&item.Stock,
			&item.Used,
			&item.OrderNumber,
			&item.BarcodeID,
			&item.InsertDate,
			&attr.Value,
			&attr.Unit,
			&attr.Label,
		); err != nil {
			return nil, fmt.Errorf("error reading items %w", err)
		}

		item.SQLParseAttributes(attr, ";;")
		items = append(items, item)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
