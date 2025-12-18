package controller

import (
	"github.com/Y0UNGMAN/DormGo/backend/logic"
	"github.com/gin-gonic/gin"

	"net/http"
)

// AIPolishRequest 请求参数
type AIPolishRequest struct {
	Content string `json:"content" binding:"required"`
}

// PolishContent 润色接口
func PolishContent(c *gin.Context) {
	var req AIPolishRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "内容不能为空"})
		return
	}

	newContent, err := logic.PolishPostContent(req.Content)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "AI 正在开小差，请稍后再试"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": newContent,
	})
}
