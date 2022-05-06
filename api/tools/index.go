package tools

import (
	"net/http"
	"os"
	"os/exec"

	"github.com/alvan/opsul/api"
	"github.com/alvan/opsul/app"
	"github.com/alvan/opsul/app/model"
	"github.com/alvan/opsul/app/utils"
	"github.com/gin-gonic/gin"
)

func Index(module string, engine *gin.Engine, router *gin.RouterGroup) {
	router.GET(module, func(ctx *gin.Context) {
		if name := ctx.Query("name"); name != "" {
			if tool := app.Store.FindToolByName(name); tool != nil {
				ctx.JSON(http.StatusOK, api.State{Data: []*model.Tool{tool}})
			} else {
				ctx.JSON(http.StatusBadRequest, api.State{Errs: []string{"Not found."}})
			}

			return
		}

		ctx.JSON(http.StatusOK, api.State{Data: app.Store.Tools})
	})

	router.GET(module+"/exec", func(ctx *gin.Context) {
		user := app.Store.FindUserByName(ctx.GetString(gin.AuthUserKey))
		if user == nil {
			ctx.JSON(http.StatusBadRequest, api.State{Errs: []string{"Invalid user."}})
			return
		}

		repo := app.Store.FindRepoByName(ctx.Query("repo"))
		if repo == nil {
			ctx.JSON(http.StatusBadRequest, api.State{Errs: []string{"Invalid repo."}})
			return
		}

		tool := app.Store.FindToolByName(ctx.Query("tool"))
		if tool == nil {
			ctx.JSON(http.StatusBadRequest, api.State{Errs: []string{"Invalid tool."}})
			return
		}

		args := ctx.QueryArray("args[]")
		if len(args) < 1 {
			if ret, err := utils.Args(ctx.Query("args")); err != nil {
				ctx.JSON(http.StatusBadRequest, api.State{Errs: []string{"Invalid args."}})
				return
			} else {
				args = ret
			}
		}

		cmd := exec.Command(tool.Path, args...)
		cmd.Dir = repo.Path
		cmd.Env = append(os.Environ(),
			"OPS_USER_NAME="+user.Name,
			"OPS_REPO_NAME="+repo.Name,
			"OPS_REPO_PATH="+repo.Path,
			"OPS_TOOL_NAME="+tool.Name,
			"OPS_TOOL_PATH="+tool.Path,
		)

		if out, err := cmd.CombinedOutput(); err != nil {
			ctx.String(http.StatusInternalServerError, string(out))
		} else {
			ctx.String(http.StatusOK, string(out))
		}
	})
}
