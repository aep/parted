package api

import (
	"log"
	"net/http"
	"strconv"

	"github.com/aep/sour"
	"github.com/gin-gonic/gin"
)

const PerPage int = 20

// GetInventory returns the inventory page
func (api *API) GetInventory(c *gin.Context) {
	query := c.Request.URL.Query()

	// parse the page
	page, _ := strconv.Atoi(query.Get("page"))
	switch {
	case page >= 1:
		page--
	case page < 1:
		page = 0
	}

	Inbounds, Max, err := api.DB.ReadInboundList(page, PerPage)
	if err != nil {
		c.JSON(500, "there was an error")
		log.Println(err)
		return
	}

	m := gin.H{
		"static":      sour.Static,
		"nav":         "inventory",
		"Inbounds":    Inbounds,
		"MaxItems":    Max,
		"FirstP":      page,
		"SecondP":     page + 1,
		"ThirdP":      page + 2,
		"HasPrevPage": true,
		"HasNextPage": true,
	}

	if page > 0 {
		m["PrevPage"] = page
	}
	if Max/PerPage >= page+2 {
		m["NextPage"] = page + 2
	}

	c.HTML(http.StatusOK, "inventory.html", m)
}

// GetJSONInventory returns the inventory in json form
func (api *API) GetJSONInventory(c *gin.Context) {
	items, err := api.DB.ReadAll()
	if err != nil {
		c.JSON(500, gin.H{"error": err})
		return
	}

	c.JSON(200, &items)
}
