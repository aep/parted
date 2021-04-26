package db

import (
	"strings"

	"github.com/aep/parted/src/elem14"
)

type Attribute struct {
	Labels string
	Units  string
	Values string
}

func SQLParseAttributes(item *elem14.Item, attr Attribute, sep string) *elem14.Item {
	values := strings.Split(attr.Values, sep)
	units := strings.Split(attr.Units, sep)
	labels := strings.Split(attr.Labels, sep)

	for i := range labels {
		item.Attributes = append(item.Attributes, elem14.Attributes{
			Attributelabel: labels[i],
			Attributevalue: values[i],
			Attributeunit:  units[i],
		})
	}
	return item
}
