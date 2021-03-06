// Package elem14 handles interactions with the element14 api
package elem14

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

// Search the API using the search term.
// configure using the passed struct, see the docs here
// https://partner.element14.com/docs/Product_Search_API_REST__Description#queryparameters
func (conf *Configuration) Search(term string) (*ManufacturerPartNumberSearch, error) {
	uVals := url.Values{}

	// mandatory
	uVals.Set("callInfo.apiKey", conf.APIKey)
	uVals.Set("term", conf.Field+":"+term)
	uVals.Set("storeInfo.id", conf.StoreInfo)
	uVals.Set("callInfo.responseDataFormat", "JSON")

	// optional
	setIfValid(uVals, "resultsSettings.numberOfResults", strconv.Itoa(conf.ResultCount))
	setIfValid(uVals, "resultsSettings.offset", strconv.Itoa(conf.Resultoffset))
	setIfValid(uVals, "callInfo.callback", conf.Callback)
	setIfValid(uVals, "userInfo.signature", conf.Signature)
	setIfValid(uVals, "userInfo.timestamp", conf.Timestamp)
	setIfValid(uVals, "userInfo.customerId", conf.CustomerID)
	setIfValid(uVals, "resultsSettings.refinements.filters", conf.Filters)
	setIfValid(uVals, "resultsSettings.responseGroup", conf.ResponseGroup)

	req, err := http.NewRequest(
		http.MethodGet,
		"https://api.element14.com/catalog/products?"+uVals.Encode(),
		nil,
	)
	if err != nil {
		return nil, err
	}

	resp, err := conf.Client.Do(req)
	if err != nil {
		return nil, err
	}

	var searchReturn ManufacturerPartNumberSearch
	if err := json.NewDecoder(resp.Body).Decode(&searchReturn); err != nil {
		return nil, err
	}

	return &searchReturn, nil
}

// setIfValid sets the parameter if it's not an empty string.
// It's used to produce clean query strings that don't have empty parameters,
// which could be misinterpreted by apis
func setIfValid(u url.Values, parameter, value string) {
	if value != "" {
		u.Set(parameter, value)
	}
}

// Configuration represent the possible parameters of the query
// see https://partner.element14.com/docs/Product_Search_API_REST__Description#queryparameters
type Configuration struct {
	// pass in your custom client, don't forget the timeout
	Client *http.Client

	// field: "any", "id", "manuPartNum"
	Field string

	// callInfo.callback
	Callback string

	// https://partner.element14.com/docs/Product_Search_API_REST__Description#storeinfo
	// storeInfo.id
	StoreInfo string

	// userInfo.signature
	Signature string

	// userInfo.timestamp
	Timestamp string

	// userInfo.billAccNum
	BillAccNum string

	// userInfo.customerId
	CustomerID string

	// userInfo.contractId
	ContractID string

	// callInfo.apiKey
	APIKey string

	// resultsSettings.offset
	Resultoffset int

	// resultsSettings.numberOfResults
	ResultCount int

	// resultsSettings.refinements.filters
	Filters string

	// resultsSettings.responseGroup
	ResponseGroup string
}

// ManufacturerPartNumberSearch returns the items searched
type ManufacturerPartNumberSearch struct {
	Manufacturerpartnumbersearchreturn ManuFacturerPartNumberSearchReturn `json:"manufacturerPartNumberSearchReturn"`
}

// Image is a component image
type Image struct {
	Basename string `json:"baseName"`
	Vrntpath string `json:"vrntPath"`
}

// Prices is the component's price
type Prices struct {
	To   int     `json:"to"`
	From int     `json:"from"`
	Cost float64 `json:"cost"`
}

// Attributes is the component's attributes
type Attributes struct {
	Attributelabel string `json:"attributeLabel"`
	Attributevalue string `json:"attributeValue"`
	Attributeunit  string `json:"attributeUnit,omitempty"`
}

// Related checks if the component has related components
type Related struct {
	Containalternatives            bool `json:"containAlternatives"`
	Containcontainrohsalternatives bool `json:"containcontainRoHSAlternatives"`
	Containaccessories             bool `json:"containAccessories"`
	Containcontainrohsaccessories  bool `json:"containcontainRoHSAccessories"`
}

// Breakdown of the component's location
type Breakdown struct {
	Inv       int    `json:"inv"`
	Region    string `json:"region"`
	Lead      int    `json:"lead"`
	Warehouse string `json:"warehouse"`
}

// Regionalbreakdown is a breakdown per region
type Regionalbreakdown struct {
	Level                       int    `json:"level"`
	Leastleadtime               int    `json:"leastLeadTime"`
	Status                      int    `json:"status"`
	Warehouse                   string `json:"warehouse"`
	Shipsfrommultiplewarehouses bool   `json:"shipsFromMultipleWarehouses"`
}

// Stock is the current item's stock
type Stock struct {
	Level                       int                 `json:"level"`
	Leastleadtime               int                 `json:"leastLeadTime"`
	Status                      int                 `json:"status"`
	Shipsfrommultiplewarehouses bool                `json:"shipsFromMultipleWarehouses"`
	Breakdown                   []Breakdown         `json:"breakdown"`
	Regionalbreakdown           []Regionalbreakdown `json:"regionalBreakdown"`
}

// Datasheets represents the item's datasheets
type Datasheets struct {
	Type        string `json:"type"`
	Description string `json:"description"`
	URL         string `json:"url"`
}

// Product represents a component
type Product struct {
	Sku                              string       `json:"sku"`
	Displayname                      string       `json:"displayName"`
	Productstatus                    string       `json:"productStatus"`
	Rohsstatuscode                   string       `json:"rohsStatusCode"`
	Packsize                         int          `json:"packSize"`
	Unitofmeasure                    string       `json:"unitOfMeasure"`
	ID                               string       `json:"id"`
	Image                            Image        `json:"image"`
	Prices                           []Prices     `json:"prices"`
	Inv                              int          `json:"inv"`
	Vendorid                         string       `json:"vendorId"`
	Vendorname                       string       `json:"vendorName"`
	Brandname                        string       `json:"brandName"`
	Translatedmanufacturerpartnumber string       `json:"translatedManufacturerPartNumber"`
	Translatedminimumorderquality    int          `json:"translatedMinimumOrderQuality"`
	Attributes                       []Attributes `json:"attributes"`
	Related                          Related      `json:"related"`
	Stock                            Stock        `json:"stock"`
	Countryoforigin                  string       `json:"countryOfOrigin"`
	Comingsoon                       bool         `json:"comingSoon"`
	Inventorycode                    int          `json:"inventoryCode"`
	Nationalclasscode                interface{}  `json:"nationalClassCode"`
	Publishingmodule                 interface{}  `json:"publishingModule"`
	Vathandlingcode                  string       `json:"vatHandlingCode"`
	Releasestatuscode                int          `json:"releaseStatusCode"`
	Isspecialorder                   bool         `json:"isSpecialOrder"`
	Isawaitingrelease                bool         `json:"isAwaitingRelease"`
	Reeling                          bool         `json:"reeling"`
	Discountreason                   int          `json:"discountReason"`
	Brandid                          string       `json:"brandId"`
	Commodityclasscode               string       `json:"commodityClassCode"`
	Datasheets                       []Datasheets `json:"datasheets,omitempty"`
}

// ManuFacturerPartNumberSearchReturn is self explainatory
type ManuFacturerPartNumberSearchReturn struct {
	Numberofresults int       `json:"numberOfResults"`
	Products        []Product `json:"products"`
}

// ToItems is used to map an inbound post to a series of items
func (e *ManufacturerPartNumberSearch) ToItems() []Item {
	items := make([]Item, 0, len(e.Manufacturerpartnumbersearchreturn.Products))
	for _, item := range e.Manufacturerpartnumbersearchreturn.Products {
		items = append(items, Item{
			Manufacturer: item.Brandname,
			PartNumber:   item.Translatedmanufacturerpartnumber,
			Description:  strings.Join(strings.Fields(item.Displayname), " "),
			Image:        item.Image.Vrntpath + item.Image.Basename,
			Stock:        item.Inv,
			BarcodeID:    item.Inventorycode,
			Attributes:   item.Attributes,
		})
	}

	return items
}

// Item represents is an item in the inventory
type Item struct {
	ID           int          `json:"id"`
	Manufacturer string       `json:"manufacturer"`
	PartNumber   string       `json:"part_number"`
	Description  string       `json:"description"`
	Image        string       `json:"image"`
	Attributes   []Attributes `json:"attributes"`
	Stock        int          `json:"stock"`
	Used         int          `json:"used"`
	OrderNumber  string       `json:"order_number"`
	BarcodeID    int          `json:"barcode_id"`
}
