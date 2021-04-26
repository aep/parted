package api

import (
	"net/http"

	"github.com/aep/sour"
	"github.com/gin-gonic/gin"
)

// GetInventory returns the inventory page
func GetInventory(c *gin.Context) {
	c.HTML(http.StatusOK, "inventory.html", gin.H{
		"static": sour.Static,
		"nav":    "inventory",
	})
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
