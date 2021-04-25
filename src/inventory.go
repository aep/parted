package src

// Item represents is an item in the inventory
type Item struct {
	Manufacturer string
	PartNb       int
	Description  string
	Image        string
	Category     string
	Spec         map[string]string
	Stock        int
	Used         int
	OrderNb      string
	Barcode      Barcode
}

// Barcode is the item's barcode
type Barcode struct {
	ID         int
	partNumber int
}

// Inventory has items
type Inventory []Item
