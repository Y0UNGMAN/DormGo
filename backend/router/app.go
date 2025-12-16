package router

import (
	"github.com/Y0UNGMAN/DormGo/backend/controller"
	"github.com/Y0UNGMAN/DormGo/backend/middleware"
	"github.com/gin-gonic/gin"
)

func App() *gin.Engine {
	r := gin.Default()
	r.Use(middleware.Cors())

	// 健康检查
	r.GET("/ping", controller.Ping)
	r.GET("/", controller.Ping)

	// ===========================
	// 用户相关路由 (User)
	// ===========================
	userPublic := r.Group("/api/v1/user")
	{
		userPublic.POST("/signup", controller.Register)
		userPublic.POST("/login", controller.Login)
	}

	userAuth := r.Group("/api/v1/user")
	userAuth.Use(middleware.JWTAuthMiddleware())
	{
		userAuth.GET("/profile", controller.GetUserProfile)
		userAuth.PUT("/profile", controller.UpdateUserProfile)
		userAuth.GET("/login-info", controller.GetLoginInfo)
		userAuth.POST("/reset-pwd", controller.ResetPassword)
	}

	// ===========================
	// 管理员相关路由 (Admin)
	// ===========================

	// 1. 管理员公共路由
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

	r.GET("/api/v1/config/types", controller.GetPostTypesAdmin)

	// ===========================
	// 帖子相关路由 (Post)
	// ===========================
	post := r.Group("/api/v1/post")
	{
		post.GET("/dorms", controller.GetDorms)
		post.GET("/post_type", controller.PostType)
		post.GET("/post_type/:id", controller.PostTypeDetial)
		post.POST("/create", controller.PostCreate)
		post.GET("/view/:id", controller.GetPostDetail)
		post.GET("/:dormid", controller.GetPostByDorm)
		post.GET("/posts", controller.GetPosts)
		post.POST("/postlike", controller.PostLike)
		post.POST("/cancellike", controller.CancelPostLike)
		post.GET("like")
		post.POST("/comment", controller.CreateComment)
		post.GET("/getcomment", controller.GetComment)
	}

	return r
}
