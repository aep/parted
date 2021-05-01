package db

import (

	// database driver initialization
	"github.com/aep/parted/src/elem14"
	_ "github.com/mattn/go-sqlite3"
)

// StoreInbound implements Storer
func (db *Database) StoreInbound(items []elem14.Item, order string) error {
	tx, err := db.DB.Begin()
	if err != nil {
		return err
	}

	//nolint
	defer tx.Rollback()

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
			&order,
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

// ReadAll the items
func (db *Database) ReadAll() ([]elem14.Item, error) {
	rows, err := db.DB.Query(SQLReadItems)
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
