package db

import (
	"strings"
)

func (item *Item) SQLParseAttributes(attr IntermediateAttr, sep string) *Item {
	values := strings.Split(attr.Value.String, sep)
	units := strings.Split(attr.Unit.String, sep)
	labels := strings.Split(attr.Label, sep)
	if len(values) < len(labels) {
		values = append(values, make([]string, len(values)-len(labels))...)
	}
	if len(units) < len(labels) {
		values = append(values, make([]string, len(values)-len(units))...)
	}

	for i := range labels {
		item.Attributes = append(item.Attributes, Attribute{
			Label: labels[i],
			Value: values[i],
			Unit:  units[i],
		})
	}
	return item
}
