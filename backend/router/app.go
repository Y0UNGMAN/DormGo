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

	user := r.Group("/api/v1/user")
	{
		//用户注册
		user.POST("/signup", controller.Register)
		//用户登录
		user.POST("/login", controller.Login)
		//获取用户信息
		user.GET("/info", controller.GetUserInfoPublic)
	}
	authPost := r.Group("/api/v1/post")
	authPost.Use(middleware.JWTAuthMiddleware())
	{
		//发帖
		authPost.POST("/create", controller.PostCreate)
		//对帖子点赞
		authPost.POST("/postlike", controller.PostLike)
		//取消点赞
		authPost.POST("/cancellike", controller.CancelPostLike)
		//发布评论
		authPost.POST("/comment", controller.CreateComment)
		//点击查看帖子详情（根据id返回帖子）
		authPost.GET("/view/:id", controller.GetPostDetail)
		//报名
		authPost.POST("/signup", controller.PostSignup)

	}

	message := r.Group("api/v1/message")
	message.Use(middleware.JWTAuthMiddleware())
	{
		//获取通知
		message.GET("/notifications", controller.GetNotifacation)
		//发私信
		message.POST("/send", controller.SendMessage)
		//获取聊天记录
		message.GET("/history", controller.GetMessageHistory)
		//获取未读通知数量
		message.GET("/unreadcount", controller.GetUnreadCount)
		//获取未读消息数量
		message.GET("/unread_count", controller.GetUnreadMessageCount)
		//已读所有消息
		message.POST("/read_all", controller.ReadAllNotification)
		//获取会话列表
		message.GET("/conversations", controller.GetConversations)
	}

	post := r.Group("/api/v1/post")
	{
		//获取所有宿舍
		post.GET("/dorms", controller.GetDorms)
		//获取帖子类型
		post.GET("/post_type", controller.PostType)
		//根据id获取帖子类型
		post.GET("/post_type/:id", controller.PostTypeDetial)

		//根据宿舍楼返回帖子
		post.GET("/:dormid", controller.GetPostByDorm)
		//获取所有帖子列表
		post.GET("/posts", controller.GetPosts)
		//获取点赞情况
		post.GET("like", controller.PostLike)
		//获取评论
		post.GET("/getcomment", controller.GetComment)
	}

	return r
}
