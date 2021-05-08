package api

import (
	"strings"
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

	items := elem14ItemToDBItem(elements.ToItems())

	for _, i := range items {
		api.Cache.Store(strings.Join(strings.Fields(i.Description), " "), &cache.Item{
			ExpiryTime: time.Now().Add(2 * time.Hour),
			Data:       i,
		})
	}

	c.JSON(200, items)
}
