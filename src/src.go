// Package src is the app source
package src

import (
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
	"github.com/joho/godotenv"
)

// Main function, runs all
func Main() {
	godotenv.Load()

	router := gin.Default()
	router.HTMLRender = ginview.Wrap(gorice.NewWithConfig(rice.MustFindBox("../views"), goview.Config{
		DisableCache: true, // TODO only for debug
		Master:       "layout.html",
	}))
	sour.StaticMount(router, "/static/", rice.MustFindBox("../static"))

	FrontEnd(router)
	BackEnd(router)

	log.Fatal(router.Run(":8080"))
}

// FrontEnd runs the front-end
func FrontEnd(r *gin.Engine) {
	r.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusFound, "/inbound")
	})

	r.GET("/inventory", func(c *gin.Context) {
		c.HTML(http.StatusOK, "inventory.html", gin.H{
			"static": sour.Static,
			"nav":    "inventory",
		})
	})

	r.GET("/inbound", func(c *gin.Context) {
		c.HTML(http.StatusOK, "inbound.html", gin.H{
			"static": sour.Static,
			"nav":    "inbound",
		})
	})
}

func BackEnd(r *gin.Engine) {
	var searchParams keywordSearchParams = keywordSearchParams{
		Client:        &http.Client{Timeout: 5 * time.Second},
		Field:         "manuPartNum",
		StoreInfo:     "uk.farnell.com",
		APIKey:        os.Getenv("API_KEY"),
		ResponseGroup: "inventory",
	}

	r.POST("/json/partsearch", func(c *gin.Context) {
		partNb := c.PostForm("search")

		elements, err := searchParams.Search(partNb)
		if err != nil {
			c.JSON(500, gin.H{
				"error": err,
			})
			return
		}

		c.JSON(200, elements)
	})
}
