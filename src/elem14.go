package src

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"
)

// Searches the API using the search term.
// configure using the passed struct, see the docs here
// https://partner.element14.com/docs/Product_Search_API_REST__Description#queryparameters
func (p *keywordSearchParams) Search(term string) (*ManufacturerPartNumberSearch, error) {
	uVals := url.Values{}

	// mandatory
	uVals.Set("callInfo.apiKey", p.APIKey)
	uVals.Set("term", p.Field+":"+term)
	uVals.Set("storeInfo.id", p.StoreInfo)
	uVals.Set("callInfo.responseDataFormat", "JSON")

	// optional
	setIfValid(uVals, "resultsSettings.numberOfResults", strconv.Itoa(p.ResultCount))
	setIfValid(uVals, "resultsSettings.offset", strconv.Itoa(p.Resultoffset))
	setIfValid(uVals, "callInfo.callback", p.Callback)
	setIfValid(uVals, "userInfo.signature", p.Signature)
	setIfValid(uVals, "userInfo.timestamp", p.Timestamp)
	setIfValid(uVals, "userInfo.customerId", p.CustomerID)
	setIfValid(uVals, "resultsSettings.refinements.filters", p.Filters)
	setIfValid(uVals, "resultsSettings.responseGroup", p.ResponseGroup)

	req, err := http.NewRequest(
		http.MethodGet,
		"https://api.element14.com/catalog/products?"+uVals.Encode(),
		nil,
	)
	if err != nil {
		return nil, err
	}

	resp, err := p.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

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
	if value != "" && value != "0" {
		u.Set(parameter, value)
	}
}

// keywordSearchParams represent the possible parameters of the query
// see https://partner.element14.com/docs/Product_Search_API_REST__Description#queryparameters
type keywordSearchParams struct {
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
	Manufacturerpartnumbersearchreturn Manufacturerpartnumbersearchreturn `json:"manufacturerPartNumberSearchReturn"`
}

type Image struct {
	Basename string `json:"baseName"`
	Vrntpath string `json:"vrntPath"`
}

type Prices struct {
	To   int     `json:"to"`
	From int     `json:"from"`
	Cost float64 `json:"cost"`
}

type Attributes struct {
	Attributelabel string `json:"attributeLabel"`
	Attributevalue string `json:"attributeValue"`
	Attributeunit  string `json:"attributeUnit,omitempty"`
}

type Related struct {
	Containalternatives            bool `json:"containAlternatives"`
	Containcontainrohsalternatives bool `json:"containcontainRoHSAlternatives"`
	Containaccessories             bool `json:"containAccessories"`
	Containcontainrohsaccessories  bool `json:"containcontainRoHSAccessories"`
}

type Breakdown struct {
	Inv       int    `json:"inv"`
	Region    string `json:"region"`
	Lead      int    `json:"lead"`
	Warehouse string `json:"warehouse"`
}

type Regionalbreakdown struct {
	Level                       int    `json:"level"`
	Leastleadtime               int    `json:"leastLeadTime"`
	Status                      int    `json:"status"`
	Warehouse                   string `json:"warehouse"`
	Shipsfrommultiplewarehouses bool   `json:"shipsFromMultipleWarehouses"`
}

type Stock struct {
	Level                       int                 `json:"level"`
	Leastleadtime               int                 `json:"leastLeadTime"`
	Status                      int                 `json:"status"`
	Shipsfrommultiplewarehouses bool                `json:"shipsFromMultipleWarehouses"`
	Breakdown                   []Breakdown         `json:"breakdown"`
	Regionalbreakdown           []Regionalbreakdown `json:"regionalBreakdown"`
}

type Datasheets struct {
	Type        string `json:"type"`
	Description string `json:"description"`
	URL         string `json:"url"`
}

type Products struct {
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

type Manufacturerpartnumbersearchreturn struct {
	Numberofresults int        `json:"numberOfResults"`
	Products        []Products `json:"products"`
}
