package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"

	apis "github.com/alvan/opsul/api/index"
	"github.com/alvan/opsul/app"
	"github.com/gin-gonic/gin"
)

var (
	conf = flag.String("conf", "etc/opsul.json", "Configuration file")

	host = flag.String("host", "", "Server host")
	port = flag.Int("port", 9900, "Server port")
)

func serv() *gin.Engine {
	engine := gin.Default()

	if app.Store.Files.Use {
		func(path string) {
			engine.StaticFile("/", path)

			list, _ := os.ReadDir(path)
			for _, file := range list {
				if !strings.HasPrefix(file.Name(), ".") {
					if file.IsDir() {
						engine.Static("/"+file.Name(), filepath.Join(path, file.Name()))
					} else {
						engine.StaticFile("/"+file.Name(), filepath.Join(path, file.Name()))
					}
				}
			}
		}(app.Store.Files.Dir)
	}

	if app.Store.Webui.Use {
		engine.SetFuncMap(template.FuncMap{
			"join": strings.Join,
			"json": json.Marshal,
		})

		engine.LoadHTMLFiles(func(path string) (list []string) {
			filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
				if filepath.Ext(path) == ".tmpl" {
					list = append(list, path)
				}
				return nil
			})

			return
		}(app.Store.Webui.Dir)...)

		engine.Group("/web", gin.BasicAuth(app.Store.AuthBasicUsers())).GET("/*path", func(ctx *gin.Context) {
			user := app.Store.FindUserByName(ctx.GetString(gin.AuthUserKey))
			path := ctx.Param("path")

			if path == "/" {
				path = "/index"
			}
			path = "/web" + path

			ctx.HTML(http.StatusOK, "/web/index", gin.H{
				"store": app.Store,
				"state": gin.H{
					"path": path,
					"user": user,
				},
			})
		})
	}

	apis.Index(engine)

	return engine
}

func main() {
	flag.Parse()

	if *conf == "" {
		flag.PrintDefaults()
		return
	}

	if err := app.Store.Load(*conf); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if err := app.Tmpfs.Init(app.Store.Tmpfs.Dir, app.Store.Tmpfs.Pre); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	sig := make(chan os.Signal, 2)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sig
		app.Tmpfs.Done()
		os.Exit(1)
	}()

	if app.Store.Https.Use {
		serv().RunTLS(fmt.Sprintf("%s:%d", *host, *port), app.Store.Https.Crt, app.Store.Https.Key)
	} else {
		serv().Run(fmt.Sprintf("%s:%d", *host, *port))
	}
}
