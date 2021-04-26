package api

import "github.com/gin-gonic/gin"

// SearchPart searches for a part
func (api *API) SearchPart(c *gin.Context) {
	partNumber := c.Param("part")

	elements, err := api.Params.Search(partNumber)
	if err != nil {
		c.JSON(404, gin.H{"error": err})
		return
	}

	c.JSON(200, elements.ToItems())
}
