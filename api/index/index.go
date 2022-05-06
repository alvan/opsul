package index

import (
	"net/http"
	"strings"

	"github.com/alvan/opsul/api"
	"github.com/alvan/opsul/app"
	"github.com/gin-gonic/gin"
)

var (
	store = make(map[string]api.Entry)
)

func index(module string, handle api.Entry) {
	store[module] = handle
}

func Index(engine *gin.Engine) {
	router := engine.Group("/api", gin.BasicAuth(app.Store.AuthBasicUsers()))
	router.GET("/", func(ctx *gin.Context) {
		var data = []gin.H{}

		for _, v := range engine.Routes() {
			if strings.HasPrefix(v.Path, router.BasePath()+"/") {
				data = append(data, gin.H{
					"path": v.Path,
					"verb": v.Method,
				})
			}
		}

		ctx.JSON(http.StatusOK, api.State{Data: data})
	})

	for module, handle := range store {
		handle(module, engine, router)
	}
}
