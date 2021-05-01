// Package api contains the routers
package api

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/aep/parted/src/cache"
	"github.com/aep/parted/src/elem14"
	"github.com/aep/sour"
	"github.com/gin-gonic/gin"
)

// ToInbound redirects the default page to the inbound page
func ToInbound(c *gin.Context) {
	c.Redirect(http.StatusFound, "/inbound")
}

// GetInbound returns the inbound html page
func GetInbound(c *gin.Context) {
	c.HTML(http.StatusOK, "inbound.html", gin.H{
		"static": sour.Static,
		"nav":    "inbound",
	})
}

// GetInboundByNumber returns the concerned inbound html page
func (api *API) GetInboundByNumber(c *gin.Context) {
	inbound := c.Param("inbound")

	items, err := api.DB.GetInboundOrder(inbound)
	if err != nil {
		c.JSON(500, gin.H{"error": err})
		return
	}

	for _, i := range items {
		api.Cache.Store(strings.Join(strings.Fields(i.Description), " "), &cache.Item{
			ExpiryTime: time.Now().Add(2 * time.Hour),
			Data:       i,
		})
	}

	c.HTML(http.StatusOK, "inbound-by-id.html", gin.H{
		"static":  sour.Static,
		"nav":     "inbound",
		"Items":   items,
		"Inbound": inbound,
	})
}

// ModifyInbound changes an existing inbound item
func (api *API) ModifyInbound(c *gin.Context) {
	inboundNbr := c.Param("inbound")

	items := c.PostFormArray("item")
	amounts := c.PostFormArray("amount")
	if len(items) != len(amounts) {
		c.JSON(400, inboundNbr)
		return
	}

	inbound := InboundPOST{
		OrderNumber: inboundNbr,
		Data:        make([]elem14.Item, 0, len(items)),
	}
	for i := range items {
		it := api.Cache.Retrieve(items[i])
		if it == nil {
			c.JSON(400, "invalid item stored")
			return
		}
		item, ok := it.Data.(elem14.Item)
		if !ok {
			c.JSON(400, "invalid item stored")
			return
		}

		var err error
		item.Stock, err = strconv.Atoi(amounts[i])
		if err != nil {
			c.JSON(400, "invalid amount field ")
			return
		}
		inbound.Data = append(inbound.Data, item)
	}

	err := api.DB.UpdateInbound(inbound.Data, inbound.OrderNumber)
	if err != nil {
		c.JSON(400, "error occured: "+err.Error())
		return
	}

	c.JSON(200, inbound)
}

// CreateInbound creates an inbound item and stores it inside the database
func (api *API) CreateInbound(c *gin.Context) {
	inbound := InboundPOST{}
	err := json.NewDecoder(c.Request.Body).Decode(&inbound)
	if err != nil {
		c.JSON(500, gin.H{"error": err})
		return
	}
	err = api.DB.Store(context.Background(), inbound.Data, inbound.OrderNumber)
	if err != nil {
		c.JSON(500, gin.H{"error": err})
		log.Println(err)
		return
	}

	c.JSON(200, inbound)
}

// InboundPOST represents an inbound post form
// It contains the order number and the products scanned
type InboundPOST struct {
	OrderNumber string        `json:"order_number"`
	Data        []elem14.Item `json:"data"`
}

func (api *API) GetInboundItem(c *gin.Context) {
	inbound := c.Param("inbound")
	items, err := api.DB.GetInboundOrder(inbound)
	if err != nil {
		c.JSON(500, gin.H{"error": err})
		return
	}

	c.JSON(200, &items)
}
