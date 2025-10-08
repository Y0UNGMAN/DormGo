package router

import (
	"github.com/Y0UNGMAN/DormGo/backend/api"
	"github.com/Y0UNGMAN/DormGo/backend/middleware"
	"github.com/gin-gonic/gin"
)

func App() *gin.Engine {
	r := gin.Default()
	r.Use(middleware.Cors())
	r.GET("/ping", api.Ping)
	r.GET("/", api.Ping)
	return r
}
