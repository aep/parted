package api

import (
	"time"

	"github.com/aep/parted/src/cache"
	"github.com/gin-gonic/gin"
)

// SearchPart searches for a part
func (api *API) SearchPart(c *gin.Context) {
	partNumber := c.Param("part")

	elements, err := api.Params.Search(partNumber)
	if err != nil {
		c.JSON(404, gin.H{"error": err})
		return
	}

	items := elements.ToItems()

	for _, i := range items {
		api.Cache.Store(i.Description, &cache.Item{
			ExpiryTime: time.Now().Add(30 * time.Minute),
			Data:       i,
		})
	}

	c.JSON(200, items)
}
