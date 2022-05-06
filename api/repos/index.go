package repos

import (
	"net/http"

	"github.com/alvan/opsul/api"
	"github.com/alvan/opsul/app"
	"github.com/alvan/opsul/app/model"
	"github.com/gin-gonic/gin"
)

func Index(module string, engine *gin.Engine, router *gin.RouterGroup) {
	router.GET(module, func(ctx *gin.Context) {
		if name := ctx.Query("name"); name != "" {
			if repo := app.Store.FindRepoByName(name); repo != nil {
				ctx.JSON(http.StatusOK, api.State{Data: []*model.Repo{repo}})
			} else {
				ctx.JSON(http.StatusBadRequest, api.State{Errs: []string{"Not found."}})
			}

			return
		}

		ctx.JSON(http.StatusOK, api.State{Data: app.Store.Repos})
	})
}
