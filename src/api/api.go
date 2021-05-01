package api

import (
	"net/http"
	"os"
	"time"

	rice "github.com/GeertJohan/go.rice"
	"github.com/aep/parted/src/cache"
	"github.com/aep/parted/src/db"
	"github.com/aep/parted/src/elem14"
	"github.com/aep/sour"
	"github.com/foolin/goview"
	"github.com/foolin/goview/supports/ginview"
	"github.com/foolin/goview/supports/gorice"
	"github.com/gin-gonic/gin"
)

// API contains the required information for the API
type API struct {
	DB     *db.Database
	Params elem14.Configuration
	Engine *gin.Engine
	Cache  *cache.Store
}

func New() *API {
	engine := gin.Default()
	engine.HTMLRender = ginview.Wrap(gorice.NewWithConfig(rice.MustFindBox("../../views"), goview.Config{
		DisableCache: true, // TODO only for debug
		Master:       "layout.html",
	}))
	sour.StaticMount(engine, "/static/", rice.MustFindBox("../../static"))

	return &API{
		Params: elem14.Configuration{
			Client:        &http.Client{Timeout: 5 * time.Second},
			Field:         "manuPartNum",
			StoreInfo:     "uk.farnell.com",
			APIKey:        os.Getenv("API_KEY"),
			ResponseGroup: "large",
		},
		DB:     db.Connect(),
		Engine: engine,
		Cache:  cache.New(5*time.Minute, func(interface{}) {}),
	}
}
