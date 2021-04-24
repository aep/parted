package src;


import (
    "github.com/joho/godotenv"
    "log"
    "github.com/gin-gonic/gin"
    "github.com/aep/sour"
    "github.com/foolin/goview"
    "github.com/foolin/goview/supports/gorice"
    "github.com/foolin/goview/supports/ginview"
    "github.com/GeertJohan/go.rice"
    "net/http"
)

func Main() {

    godotenv.Load()

    router := gin.Default()
    router.HTMLRender = ginview.Wrap(gorice.NewWithConfig(rice.MustFindBox("../views"), goview.Config {
        DisableCache:   true, //TODO only for debug
        Master:         "layout.html",
    }));
    sour.StaticMount(router, "/static/", rice.MustFindBox("../static"))

    Frontend(router)

    log.Fatal(router.Run(":8080"))
}


func Frontend(r *gin.Engine) {
    r.GET("/", func(c *gin.Context) {
        c.Redirect(http.StatusFound, "/inbound");
    })
    r.GET("/inventory", func(c *gin.Context) {
        c.HTML(http.StatusOK, "inventory.html", gin.H{
            "static":       sour.Static,
            "nav":          "inventory",
        })
    })
    r.GET("/inbound", func(c *gin.Context) {
        c.HTML(http.StatusOK, "inbound.html", gin.H{
            "static":       sour.Static,
            "nav":          "inbound",
        })
    })

    r.GET("/json/partsearch", func(c *gin.Context) {
        partnr := gin.PostForm("search")
        log.Println("PART LOOKUP");
        c.JSON(200, gin.H{
            "manufacturer": "Big Corp.",

        })
    })
}
