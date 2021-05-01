package db

import (
	"github.com/aep/parted/src/elem14"
)

// GetInboundOrder returns the data from a single order
func (db *Database) GetInboundOrder(orderNumber string) ([]elem14.Item, error) {
	rows, err := db.DB.Query(SQLReadInboud, orderNumber)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []elem14.Item
	for rows.Next() {
		var item elem14.Item
		var attr Attribute
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
			&attr.Values,
			&attr.Units,
			&attr.Labels,
		); err != nil {
			return nil, err
		}

		items = append(items, *SQLParseAttributes(&item, attr, ";;"))
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
		return err
	}

	insertItemStmt, err := tx.Prepare(SQLInsertItem)
	if err != nil {
		return err
	}

	insertMetaStmt, err := tx.Prepare(SQLInsertAttr)
	if err != nil {
		return err
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
