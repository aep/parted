package db

import (
	"fmt"
	"time"

	"github.com/aep/parted/src/elem14"
)

// GetInboundOrder returns the data from a single order
func (db *Database) GetInboundOrder(orderNumber string) ([]Item, error) {
	rows, err := db.DB.Query(SQLReadInboud, orderNumber)
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

		items = append(items, *item.SQLParseAttributes(attr, ";;"))
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

// UpdateInbound delete existing items and insert the new ones
func (db *Database) UpdateInbound(items []elem14.Item, orderNumber string) error {
	tx, err := db.DB.Begin()
	if err != nil {
		return err
	}

	defer tx.Rollback()

	if _, err := tx.Exec(SQLDeleteInbound, orderNumber); err != nil {
		return fmt.Errorf("error deleting items %w", err)
	}

	insertItemStmt, err := tx.Prepare(SQLInsertItem)
	if err != nil {
		return fmt.Errorf("error preparing insert items %w", err)
	}

	insertMetaStmt, err := tx.Prepare(SQLInsertAttr)
	if err != nil {
		return fmt.Errorf("error preparing insert attr items %w", err)
	}

	for _, item := range items {
		exec, err := insertItemStmt.Exec(
			&item.Manufacturer,
			&item.PartNumber,
			&item.Description,
			&item.Image,
			&item.Stock,
			&orderNumber,
			&item.BarcodeID,
			time.Now().UTC(),
		)
		if err != nil {
			return fmt.Errorf("error inserting item %w", err)
		}

		elementID, err := exec.LastInsertId()
		if err != nil {
			return fmt.Errorf("error getting last inserted id %w", err)
		}

		for _, spec := range item.Attributes {
			_, err := insertMetaStmt.Exec(
				&elementID,
				&spec.Attributelabel,
				&spec.Attributevalue,
				&spec.Attributeunit,
			)
			if err != nil {
				return fmt.Errorf("error getting attributes %w", err)
			}
		}

	}

	return tx.Commit()
}
