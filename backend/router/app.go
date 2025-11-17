package router

import (
	"github.com/Y0UNGMAN/DormGo/backend/controller"
	"github.com/Y0UNGMAN/DormGo/backend/middleware"
	"github.com/gin-gonic/gin"
)

func App() *gin.Engine {
	r := gin.Default()
	r.Use(middleware.Cors())
	r.GET("/ping", controller.Ping)
	r.GET("/", controller.Ping)

	post := r.Group("/api/v1/post")
	{
		//获取帖子类型
		post.GET("/post_type", controller.PostType)
		post.GET("post_type/:id", controller.PostTypeDetial)
		post.POST("create", controller.PostCreate)
	}

	return r
}
