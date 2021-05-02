package api

import (
	"time"

	"github.com/aep/parted/src/db"
	"github.com/aep/parted/src/elem14"
)

// elem14ItemToDBItem maps items requested from the elem14 api to a database item
func elem14ItemToDBItem(src []elem14.Item) []db.Item {
	dst := make([]db.Item, 0, len(src))
	for _, i := range src {
		attrDst := make([]db.Attribute, 0, len(i.Attributes))
		for _, srcAttr := range i.Attributes {
			attrDst = append(attrDst, db.Attribute{
				Label: srcAttr.Attributelabel,
				Value: srcAttr.Attributevalue,
				Unit:  srcAttr.Attributeunit,
			})
		}
		dst = append(dst, db.Item{
			ID:           i.ID,
			Manufacturer: i.Manufacturer,
			PartNumber:   i.PartNumber,
			Description:  i.Description,
			Image:        i.Image,
			Stock:        i.Stock,
			Used:         i.Used,
			OrderNumber:  i.OrderNumber,
			BarcodeID:    i.BarcodeID,
			InsertDate:   time.Now(),
			Attributes:   attrDst,
		})
	}

	return dst
}
