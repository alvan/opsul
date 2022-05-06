package props

import (
	"net/http"

	"github.com/alvan/opsul/api"
	"github.com/alvan/opsul/app"
	"github.com/gin-gonic/gin"
)

func Index(module string, engine *gin.Engine, router *gin.RouterGroup) {
	engine.GET(router.BasePath()+module, func(ctx *gin.Context) {
		if name := ctx.Query("name"); name != "" {
			if value, ok := app.Store.Props[name]; ok {
				ctx.JSON(http.StatusOK, api.State{Data: gin.H{name: value}})
			} else {
				ctx.JSON(http.StatusNotFound, api.State{Errs: []string{"Not found."}})
			}

			return
		}

		ctx.JSON(http.StatusOK, api.State{Data: app.Store.Props})
	})
}
