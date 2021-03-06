// Package api contains the routers
package api

import (
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/aep/parted/src/cache"
	"github.com/aep/parted/src/db"
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

// DeleteInbound deletes the requested inbound
func (api *API) DeleteInbound(c *gin.Context) {
	inbound := c.Param("inbound")

	if err := api.DB.DeleteInbound(inbound); err != nil {
		c.JSON(500, "an unexpected error happened: "+err.Error())
	}

	c.Status(200)
}

// GetInboundByNumber returns the concerned inbound html page
func (api *API) GetInboundByNumber(c *gin.Context) {
	inbound := c.Param("inbound")

	items, err := api.DB.ReadInbound(inbound)
	if err != nil {
		log.Println(err)
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
		Data:        make([]db.Item, 0, len(items)),
	}
	for i := range items {
		it := api.Cache.Retrieve(items[i])
		if it == nil {
			if err := api.RefreshInboundCache(inboundNbr); err != nil {
				it = api.Cache.Retrieve(items[i])
				if it == nil {
					c.HTML(500, "Could not retrieve the item from the cache, please re-enter your data", nil)
				}
			}
		}

		var err error
		it.Data.Stock, err = strconv.Atoi(amounts[i])
		if err != nil {
			c.JSON(400, "invalid amount field")
			log.Println(err)
			return
		}
		inbound.Data = append(inbound.Data, it.Data)
	}

	err := api.DB.UpdateInbound(inbound.Data, inbound.OrderNumber)
	if err != nil {
		c.JSON(400, "error occured: "+err.Error())
		log.Println(err)
		return
	}

	c.JSON(200, inbound)
}

func (api *API) RefreshInboundCache(inboundNbr string) error {
	items, err := api.DB.ReadInbound(inboundNbr)
	if err != nil {
		return err
	}

	for _, i := range items {
		api.Cache.Store(strings.Join(strings.Fields(i.Description), " "), &cache.Item{
			ExpiryTime: time.Now().Add(2 * time.Hour),
			Data:       i,
		})
	}
	return nil
}

// CreateInbound creates an inbound item and stores it inside the database
func (api *API) CreateInbound(c *gin.Context) {
	inboundNbr := c.PostForm("ordernumber")
	items := c.PostFormArray("item")
	amounts := c.PostFormArray("amount")
	if len(items) != len(amounts) {
		c.JSON(400, inboundNbr)
		return
	}

	inbound := InboundPOST{
		OrderNumber: inboundNbr,
		Data:        make([]db.Item, 0, len(items)),
	}
	for i := range items {
		it := api.Cache.Retrieve(strings.Join(strings.Fields(items[i]), " "))
		if it == nil {
			return
		}

		var err error
		it.Data.Stock, err = strconv.Atoi(amounts[i])
		if err != nil {
			c.JSON(400, "invalid amount field")
			log.Println(err)
			return
		}
		inbound.Data = append(inbound.Data, it.Data)
	}

	err := api.DB.InsertInbound(inbound.Data, inboundNbr)
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
	OrderNumber string    `json:"order_number"`
	Data        []db.Item `json:"data"`
}

func (api *API) GetInboundItem(c *gin.Context) {
	inbound := c.Param("inbound")
	items, err := api.DB.ReadInbound(inbound)
	if err != nil {
		c.JSON(500, gin.H{"error": err})
		log.Println(err)
		return
	}

	c.JSON(200, &items)
}
