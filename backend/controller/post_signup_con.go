package controller

import (
	"net/http"

	"github.com/Y0UNGMAN/DormGo/backend/middleware" // 引入 middleware 包以使用 Key 常量
	"github.com/Y0UNGMAN/DormGo/backend/model"
	"github.com/gin-gonic/gin"
)

func PostSignup(c *gin.Context) {
	// 1. 定义请求参数结构
	var req struct {
		PostID uint `json:"post_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "参数错误"})
		return
	}

	// 2. 安全地获取 UserID
	var userID uint
	// 使用常量 middleware.CtxUserIDKey 保持统一
	if v, exists := c.Get(middleware.CtxUserIDKey); exists {
		// 使用 comma-ok 断言，防止 Panic
		if uid, ok := v.(uint); ok {
			userID = uid
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "msg": "用户身份类型异常"})
			return
		}
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "msg": "请先登录"})
		return
	}

	// 3. 调用业务逻辑
	err := model.SignupPost(req.PostID, userID)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "报名成功"})
}
