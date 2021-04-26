// Package src is the app source
package src

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	rice "github.com/GeertJohan/go.rice"
	"github.com/aep/sour"
	"github.com/foolin/goview"
	"github.com/foolin/goview/supports/ginview"
	"github.com/foolin/goview/supports/gorice"
	"github.com/gin-gonic/gin"
)

// Main function, runs all
func Main() {
	engine := gin.Default()
	engine.HTMLRender = ginview.Wrap(gorice.NewWithConfig(rice.MustFindBox("../views"), goview.Config{
		DisableCache: true, // TODO only for debug
		Master:       "layout.html",
	}))
	sour.StaticMount(engine, "/static/", rice.MustFindBox("../static"))

	api := API{
		Params: keywordSearchParams{
			Client:        &http.Client{Timeout: 5 * time.Second},
			Field:         "manuPartNum",
			StoreInfo:     "uk.farnell.com",
			APIKey:        os.Getenv("API_KEY"),
			ResponseGroup: "large",
		},
		DB:     Connect(),
		Engine: engine,
	}

	api.FrontEnd()
	api.BackEnd()

	log.Fatal(engine.Run(":8080"))
}

type API struct {
	DB     *Database
	Params keywordSearchParams
	Engine *gin.Engine
}

// FrontEnd runs the front-end
func (api *API) FrontEnd() {
	api.Engine.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusFound, "/inbound")
	})

	api.Engine.GET("/inventory", func(c *gin.Context) {
		c.HTML(http.StatusOK, "inventory.html", gin.H{
			"static": sour.Static,
			"nav":    "inventory",
		})
	})

	api.Engine.GET("/inbound", func(c *gin.Context) {
		c.HTML(http.StatusOK, "inbound.html", gin.H{
			"static": sour.Static,
			"nav":    "inbound",
		})
	})
}

// BackEnd initializes and runs the back end
func (api *API) BackEnd() {
	api.Engine.GET("/json/partsearch/:part", func(c *gin.Context) {
		partNumber := c.Param("part")

		elements, err := api.Params.Search(partNumber)
		if err != nil {
			c.JSON(404, gin.H{"error": err})
			return
		}

		c.JSON(200, elements.toItems())
	})

	api.Engine.POST("/inbound", func(c *gin.Context) {
		inbound := InboundPOST{}
		err := json.NewDecoder(c.Request.Body).Decode(&inbound)
		if err != nil {
			c.JSON(500, gin.H{"error": err})
			return
		}
		err = api.DB.Store(context.TODO(), inbound)
		if err != nil {
			c.JSON(500, gin.H{"error": err})
			log.Println(err)
			return
		}

		c.JSON(200, inbound)
	})

	api.Engine.GET("/json/inventory", func(c *gin.Context) {
		items, err := api.DB.ReadAll()
		if err != nil {
			c.JSON(500, gin.H{"error": err})
			return
		}

		c.JSON(200, &items)
	})
}

// InboundPOST represents an inbound post form
// It contains the order number and the products scanned
type InboundPOST struct {
	OrderNumber string `json:"order_number"`
	Data        []Item `json:"data"`
}
