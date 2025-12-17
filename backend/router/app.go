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
	}

	userAuth := r.Group("/api/v1/user")
	userAuth.Use(middleware.JWTAuthMiddleware())
	{
		userAuth.GET("/profile", controller.GetUserProfile)
		userAuth.PUT("/profile", controller.UpdateUserProfile)
		userAuth.GET("/login-info", controller.GetLoginInfo)
		userAuth.POST("/reset-pwd", controller.ResetPassword)
		userAuth.GET("/info", controller.GetUserInfoPublic)
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
		//收藏
		authPost.POST("/favorite", controller.PostFavorite)
		//取消收藏
		authPost.POST("/cancel_favorite", controller.CancelPostFavorite)
		// 获取收藏列表
		authPost.GET("/my_favorites", controller.GetMyFavorites)
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

	adminPublic := r.Group("/api/v1/admin")
	{
		adminPublic.POST("/login", controller.AdminLogin)
	}

	// 2. 管理员鉴权路由
	adminAuth := r.Group("/api/v1/admin")
	adminAuth.Use(middleware.JWTAuthMiddleware())
	{
		// 个人中心 (修复获取资料错误问题)
		adminAuth.GET("/profile", controller.GetAdminProfile)
		adminAuth.PUT("/profile", controller.UpdateAdminProfile)

		// 数据统计
		adminAuth.GET("/statistics", controller.AdminStatistics)

		// 用户管理
		adminAuth.GET("/users", controller.GetUserList)
		adminAuth.PUT("/users/:id/status", controller.UpdateUserStatus)
		adminAuth.PUT("/users/batch/status", controller.BatchUpdateUserStatus)
		adminAuth.PUT("/users/:id/credit", controller.UpdateUserCredit)

		// 楼栋管理
		adminAuth.GET("/dorms", controller.GetDormListAdmin)
		adminAuth.POST("/dorms", controller.CreateDorm)
		adminAuth.PUT("/dorms/:id", controller.UpdateDorm)
		adminAuth.DELETE("/dorms/:id", controller.DeleteDorm)

		// 内容审核 / 帖子管理
		adminAuth.GET("/contents", controller.GetContentList)
		adminAuth.PUT("/contents/:id/audit", controller.AuditContent)
		adminAuth.PUT("/contents/batch/audit", controller.BatchAuditContent)
		// 新增：置顶/取消置顶
		adminAuth.PUT("/contents/:id/pin", controller.TogglePostPin)
		// 新增：管理员发帖
		adminAuth.POST("/contents/create", controller.AdminCreatePost)

		// 通知管理
		adminAuth.GET("/notices", controller.GetNoticeList)
		adminAuth.POST("/notices", controller.CreateNotice)
		adminAuth.DELETE("/notices/:id", controller.DeleteNotice)

		// 违规处理
		adminAuth.GET("/violations", controller.GetViolationList)
		adminAuth.PUT("/violations/:id/process", controller.ProcessViolation)

		// 系统配置
		adminAuth.PUT("/config/basic", controller.UpdateBasicConfig)
		adminAuth.PUT("/config/security", controller.UpdateSecurityConfig)

		// 敏感词
		adminAuth.GET("/config/sensitive-words", controller.GetSensitiveWords)
		adminAuth.POST("/config/sensitive-words", controller.AddSensitiveWord)
		adminAuth.PUT("/config/sensitive-words/:id", controller.UpdateSensitiveWord)
		adminAuth.DELETE("/config/sensitive-words/:id", controller.DeleteSensitiveWord)

		// 互助类型
		adminAuth.POST("/config/types", controller.CreatePostType)
		adminAuth.PUT("/config/types/:id", controller.UpdatePostType)
		adminAuth.PUT("/config/types/:id/status", controller.UpdatePostTypeStatus)
		adminAuth.DELETE("/config/types/:id", controller.DeletePostType)
	}
	return r
}
