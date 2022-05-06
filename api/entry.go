package api

import "github.com/gin-gonic/gin"

type Entry func(module string, engine *gin.Engine, router *gin.RouterGroup)
