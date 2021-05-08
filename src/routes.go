// Package src is the app source
package src

import (
	"log"

	"github.com/aep/parted/src/api"
)

// ListenAndServe function, runs all.
// returns a blocking function
func ListenAndServe() {
	a := api.New()

	a.Engine.GET("/", api.ToInbound)
	a.Engine.GET("/inventory", a.GetInventory)
	a.Engine.GET("/inbound", api.GetInbound)
	a.Engine.GET("/inbound/:inbound", a.GetInboundByNumber)
	a.Engine.POST("/inbound/:inbound", a.ModifyInbound)
	a.Engine.POST("/inbound", a.CreateInbound)
	a.Engine.GET("/json/partsearch/:part", a.SearchPart)
	a.Engine.GET("/json/inventory", a.GetJSONInventory)
	a.Engine.GET("/json/inbound/:inbound", a.GetInboundItem)

	if err := a.Engine.Run(); err != nil {
		log.Fatalln(err)
	}
}
