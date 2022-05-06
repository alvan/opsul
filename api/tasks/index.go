package tasks

import (
	"bytes"
	"crypto/hmac"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/alvan/opsul/api"
	"github.com/alvan/opsul/app"
	"github.com/alvan/opsul/app/model"
	"github.com/alvan/opsul/app/utils"
	"github.com/gin-gonic/gin"
)

func Index(module string, engine *gin.Engine, router *gin.RouterGroup) {
	router.GET(module, func(ctx *gin.Context) {
		if name := ctx.Query("name"); name != "" {
			if task := app.Store.FindTaskByName(name); task != nil &&
				(ctx.Query("id") == "" || task.Id == ctx.Query("id")) {

				ctx.JSON(http.StatusOK, api.State{Data: []*model.Task{task}})
			} else {
				ctx.JSON(http.StatusBadRequest, api.State{Errs: []string{"Not found."}})
			}

			return
		}

		ctx.JSON(http.StatusOK, api.State{Data: app.Store.Tasks})
	})

	router.POST(module, func(ctx *gin.Context) {
		body, _ := io.ReadAll(ctx.Request.Body)
		ctx.Request.Body = io.NopCloser(bytes.NewBuffer(body))

		user := app.Store.FindUserByName(ctx.GetString(gin.AuthUserKey))
		if user == nil {
			ctx.JSON(http.StatusBadRequest, api.State{Errs: []string{"Invalid user."}})
			return
		}

		sign := ctx.Request.Header.Get("X-Ops-Signature")
		if sign == "" {
			sign = ctx.Request.Header.Get("X-Hub-Signature")
		}
		if (sign == "" && user.Code != "") || (sign != "" && !hmac.Equal([]byte(sign), []byte(utils.Sign(sign, body, user.Code)))) {
			ctx.JSON(http.StatusBadRequest, api.State{Errs: []string{"Invalid sign."}})
			return
		}

		name := ctx.Query("name")
		if name == "" {
			ctx.JSON(http.StatusBadRequest, api.State{Errs: []string{"Invalid name."}})
			return
		}

		repo := app.Store.FindRepoByName(ctx.Query("repo"))
		if repo == nil {
			ctx.JSON(http.StatusBadRequest, api.State{Errs: []string{"Invalid repo."}})
			return
		}

		pack := repo.FindPackByName(ctx.Query("pack"))
		if pack == nil {
			ctx.JSON(http.StatusBadRequest, api.State{Errs: []string{"Invalid pack."}})
			return
		}

		tags := ctx.QueryArray("tags[]")
		if len(tags) < 1 {
			tags = utils.Tags(ctx.Query("tags"), " ")
		} else {
			tags = utils.Tags(strings.Join(tags, " "), " ")
		}

		hook := ctx.Query("hook")
		if hook != "" {
			if _, err := url.ParseRequestURI(hook); err != nil {
				ctx.JSON(http.StatusBadRequest, api.State{Errs: []string{"Invalid hook."}})
				return
			}
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

		otid := ""
		if task := app.Store.FindTaskByName(name); task != nil {
			otid = task.Id

			if task.Stat == model.TASK_STAT_RUNNING || task.Stat == model.TASK_STAT_PENDING {
				ctx.JSON(http.StatusConflict, api.State{Errs: []string{"Task already exists."}})
				return
			} else if task.File != "" {
				os.Remove(task.File)
			}
		}

		task := &model.Task{
			Id:   utils.Uqid.Next(),
			Name: name,
			Repo: repo.Name,
			Pack: pack.Name,
			User: user.Name,
			Tags: tags,
			Stat: model.TASK_STAT_PENDING,
			Hook: hook,
			Time: time.Now(),
			Args: args,
			Logs: []model.Note{
				model.Note{
					Name: "task.new",
					Time: time.Now(),
				},
			},
		}

		if !app.Store.SaveTask(task, otid) {
			ctx.JSON(http.StatusConflict, api.State{Errs: []string{"Failed to save new task."}})
			return
		}

		envs := []string{
			"OPS_USER_NAME=" + user.Name,
			"OPS_REPO_NAME=" + repo.Name,
			"OPS_REPO_PATH=" + repo.Path,
			"OPS_TASK_NAME=" + task.Name,
			"OPS_TASK_TAGS=" + strings.Join(task.Tags, " "),
			"OPS_PACK_NAME=" + pack.Name,
			"OPS_PACK_PATH=" + pack.Path,
			"OPS_DATA_SIGN=" + sign,
			"OPS_DATA_TYPE=" + ctx.Request.Header.Get("Content-Type"),
		}

		go func() {
			if task.Stat != model.TASK_STAT_PENDING {
				return
			}

			task.Stat = model.TASK_STAT_RUNNING
			task.Logs = append(task.Logs, model.Note{
				Name: "task.run",
				Time: time.Now(),
			})

			file, err := app.Tmpfs.File()
			if err == nil {
				defer file.Close()
				task.File = file.Name()

				cmd := exec.Command(pack.Path, task.Args...)
				cmd.Dir = repo.Path
				cmd.Env = append(os.Environ(), envs...)
				cmd.Stdin = strings.NewReader(string(body))
				cmd.Stdout = file
				cmd.Stderr = file
				cmd.SysProcAttr = utils.ProcAttrGroup()
				err = cmd.Start()
				if err == nil {
					task.Proc = cmd.Process.Pid
					err = cmd.Wait()
				}
			}

			if err == nil {
				task.Stat = model.TASK_STAT_SUCCESS
			} else {
				task.Stat = model.TASK_STAT_FAILURE
				task.Logs = append(task.Logs, model.Note{
					Name: "task.err",
					Text: err.Error(),
					Time: time.Now(),
				})
			}

			task.Logs = append(task.Logs, model.Note{
				Name: "task.end",
				Time: time.Now(),
			})

			if task.Hook != "" {
				if dat, err := json.Marshal(task); err == nil {
					sig := utils.Sign("sha1=", dat, user.Code)
					if req, err := http.NewRequest("POST", task.Hook, bytes.NewBuffer(dat)); err == nil {
						req.Header.Add("X-Ops-Signature", sig)
						req.Header.Add("Content-Type", "application/json")
						(&http.Client{Timeout: 3 * time.Second}).Do(req)
					}
				}
			}
		}()

		ctx.JSON(http.StatusCreated, api.State{Data: task})
	})

	router.GET(module+"/read", func(ctx *gin.Context) {
		name := ctx.Query("name")
		if task := app.Store.FindTaskByName(name); task != nil &&
			(ctx.Query("id") == "" || task.Id == ctx.Query("id")) &&
			task.File != "" {

			auto := ctx.Query("auto")
			if auto != "" {
				if val, err := strconv.Atoi(auto); err == nil && val > 0 {
					auto = strconv.Itoa(val)
				} else {
					auto = ""
				}
			}

			if ctx.Query("mode") == "tail" {
				if size, err := strconv.Atoi(ctx.DefaultQuery("size", "10")); err == nil && size > 0 {
					if out, err := utils.Tail(task.File, size); err == nil {
						if auto != "" && (task.Stat == model.TASK_STAT_RUNNING || task.Stat == model.TASK_STAT_PENDING) {
							ctx.Header("Refresh", auto)
						}

						ctx.String(http.StatusOK, string(out))
					} else {
						ctx.JSON(http.StatusInternalServerError, api.State{Errs: []string{err.Error()}})
					}
				} else {
					ctx.JSON(http.StatusBadRequest, api.State{Errs: []string{"Invaild size."}})
				}
			} else {
				// The default mode
				if out, err := os.ReadFile(task.File); err == nil {
					if auto != "" && (task.Stat == model.TASK_STAT_RUNNING || task.Stat == model.TASK_STAT_PENDING) {
						ctx.Header("Refresh", auto)
					}

					ctx.String(http.StatusOK, string(out))
				} else {
					ctx.JSON(http.StatusInternalServerError, api.State{Errs: []string{err.Error()}})
				}
			}
		} else {
			ctx.JSON(http.StatusBadRequest, api.State{Errs: []string{"Failed to read data."}})
		}
	})

	router.POST(module+"/stop", func(ctx *gin.Context) {
		name := ctx.Query("name")
		if task := app.Store.FindTaskByName(name); task != nil &&
			(ctx.Query("id") == "" || task.Id == ctx.Query("id")) &&
			task.Stat == model.TASK_STAT_RUNNING {

			if proc := task.Proc; proc > 0 {
				if err := utils.ProcKill(task.Proc); err == nil {
					task.Proc = 0
					ctx.JSON(http.StatusOK, api.State{Data: task})
				} else {
					ctx.JSON(http.StatusInternalServerError, api.State{Errs: []string{err.Error()}})
				}
			}

			return
		}

		ctx.JSON(http.StatusBadRequest, api.State{Errs: []string{"Failed to stop task."}})
	})

	router.POST(module+"/drop", func(ctx *gin.Context) {
		name := ctx.Query("name")
		if task := app.Store.FindTaskByName(name); task != nil &&
			(ctx.Query("id") == "" || task.Id == ctx.Query("id")) &&
			task.Stat != model.TASK_STAT_PENDING && task.Stat != model.TASK_STAT_RUNNING {

			if task.File != "" {
				os.Remove(task.File)
			}
			if app.Store.DropTask(task) {
				ctx.JSON(http.StatusOK, api.State{Data: task})
				return
			}
		}

		ctx.JSON(http.StatusBadRequest, api.State{Errs: []string{"Failed to drop task."}})
	})
}
