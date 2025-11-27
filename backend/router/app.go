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
		//获取所有宿舍
		post.GET("/dorms", controller.GetDorms)
		//获取帖子类型
		post.GET("/post_type", controller.PostType)
		//根据id获取帖子类型
		post.GET("/post_type/:id", controller.PostTypeDetial)
		//发帖
		post.POST("/create", controller.PostCreate)
		//点击查看帖子详情（根据id返回帖子）
		post.GET("/view/:id", controller.GetPostDetail)
		//根据宿舍楼返回帖子
		post.GET("/:dormid", controller.GetPostByDorm)
		//获取所有帖子列表
		post.GET("/posts", controller.GetPosts)
		//对帖子点赞
		post.POST("/postlike", controller.PostLike)
		//取消点赞
		post.POST("/cancellike", controller.CancelPostLike)
		//获取点赞情况
		post.GET("like")
	}

	return r
}
