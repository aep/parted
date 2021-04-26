// Package api contains the routers
package api

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

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
