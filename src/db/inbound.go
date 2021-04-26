package db

import "github.com/aep/parted/src/elem14"

// GetInboundOrder returns the data from a single order
func (db *Database) GetInboundOrder(orderNumber string) ([]elem14.Item, error) {
	rows, err := db.DB.Query(SQLReadItems, orderNumber)
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
