package controller

import (
	"net/http"
	"strconv"

	"github.com/Y0UNGMAN/DormGo/backend/logic"
	"github.com/gin-gonic/gin"
)

// DeletePost 删除帖子接口
func DeletePost(c *gin.Context) {
	// 1. 获取 URL 参数中的 id
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "参数错误"})
		return
	}

	// 2. 获取当前登录用户 ID
	userID, err := GetCurrentUserId(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "msg": "未登录"})
		return
	}

	// 3. 调用逻辑层
	err = logic.DeletePost(uint(id), userID)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "msg": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "删除成功",
	})
}
