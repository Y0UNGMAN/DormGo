package controller

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/Y0UNGMAN/DormGo/backend/logic"
	"github.com/gin-gonic/gin"
)

// FavoriteRequest 定义请求参数结构体
type FavoriteRequest struct {
	PostID uint `json:"post_id" binding:"required"`
}

// PostFavorite 收藏帖子
func PostFavorite(c *gin.Context) {
	// 1. 参数绑定与校验
	var req FavoriteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "参数错误", "error": err.Error()})
		return
	}

	// 2. 获取当前用户ID (从中间件设置的Context中)
	userID, err := GetCurrentUserId(c)
	if err != nil {
		fmt.Println(err)
		return
	}
	// 3. 调用 Logic 层
	if err := logic.DoFavorite(userID, req.PostID); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": err.Error()}) // 业务失败通常返回200+错误码，或者500
		return
	}

	// 4. 返回响应
	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "收藏成功"})
}

// CancelPostFavorite 取消收藏
func CancelPostFavorite(c *gin.Context) {
	var req FavoriteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "参数错误"})
		return
	}

	userID, err := GetCurrentUserId(c)
	if err != nil {
		fmt.Println(err)
		return
	}

	if err := logic.CancelFavorite(userID, req.PostID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "操作失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "取消收藏成功"})
}

// GetMyFavorites 获取我的收藏列表
func GetMyFavorites(c *gin.Context) {
	userID, err := GetCurrentUserId(c)
	if err != nil {
		fmt.Println(err)
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))

	// 调用 Logic 层
	data, err := logic.GetUserFavoriteList(userID, page, size)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "获取列表失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":  200,
		"msg":   "success",
		"data":  data["list"],
		"total": data["total"],
	})
}
