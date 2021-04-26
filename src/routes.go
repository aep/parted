// Package src is the app source
package src

import (
	"log"

	"github.com/aep/parted/src/api"
)

// Main function, runs all
func Main() {
	a := api.New()

	// Front
	a.Engine.GET("/", api.ToInbound)
	a.Engine.GET("/inventory", api.GetInventory)
	a.Engine.GET("/inbound", api.GetInbound)

	// Back
	a.Engine.POST("/inbound", a.CreateInbound)
	a.Engine.GET("/json/partsearch/:part", a.SearchPart)
	a.Engine.GET("/json/inventory", a.GetJSONInventory)
	a.Engine.GET("/json/inbound/:inbound", a.GetInboundItem)

	if err := a.Engine.Run(":8080"); err != nil {
		log.Fatalln(err)
	}
}
