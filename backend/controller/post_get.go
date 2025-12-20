package controller

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/Y0UNGMAN/DormGo/backend/logic"
	"github.com/Y0UNGMAN/DormGo/backend/middleware"
	"github.com/gin-gonic/gin"
)

// 获取帖子详情
func GetPostDetail(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		fmt.Println(err)
		c.JSON(400, gin.H{
			"code":    "400",
			"message": "Atoi failed",
		})
		return
	}
	var userId uint = 0
	if v, exists := c.Get(middleware.CtxUserIDKey); exists {
		if uid, ok := v.(uint); ok {
			userId = uid
		} else if uidInt, ok := v.(int); ok {
			userId = uint(uidInt)
		} else {
			fmt.Printf("UserID type mismatch in Context. Got: %T\n", v)
		}
	} else {
		fmt.Println("UserID not found in context")
	}
	data, isliked, isFavorited, err := logic.GetPostDetail(id, userId)
	if err != nil {
		fmt.Println(err)
		c.JSON(400, gin.H{
			"code":    "400",
			"message": "get post detail failed",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code":        http.StatusOK,
		"message":     "get post detail success",
		"data":        data,
		"isliked":     isliked,
		"isFavorited": isFavorited,
	})

}

// 根据宿舍号获取帖子
func GetPostByDorm(c *gin.Context) {
	idStr := c.Param("dormid")
	dormid, err := strconv.Atoi(idStr)
	if err != nil {
		fmt.Println(err)
		c.JSON(400, gin.H{
			"code":    "400",
			"message": "Atoi failed",
		})
		return
	}
	data, err := logic.GetPostByDorm(dormid)
	if err != nil {
		fmt.Println(err)
		c.JSON(400, gin.H{
			"code":    "400",
			"message": "get post by dorm failed",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "get post by dorm success",
		"data":    data,
	})
}

// 获取所有帖子
func GetPosts(c *gin.Context) {
	pageStr := c.DefaultQuery("page", "1")
	pageSizeStr := c.DefaultQuery("size", "10")
	page, _ := strconv.Atoi(pageStr)
	pageSize, _ := strconv.Atoi(pageSizeStr)
	if pageSize > 100 {
		pageSize = 100
	}
	data, total, err := logic.GetPosts(page, pageSize)
	if err != nil {
		fmt.Println(err)
		c.JSON(400, gin.H{
			"code":    "400",
			"message": "get posts failed",
		})
		return
	}
	hasMore := false
	if int64(page*pageSize) < total {
		hasMore = true
	}
	c.JSON(http.StatusOK, gin.H{
		"code":     http.StatusOK,
		"message":  "get posts success",
		"data":     data,
		"total":    total,   // 总数
		"page":     page,    // 当前页码
		"has_more": hasMore, // 是否还有下一页
	})
}

// GetUserPosts 获取指定用户的帖子
func GetUserPosts(c *gin.Context) {
	userIDStr := c.Query("user_id")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "参数错误"})
		return
	}

	data, err := logic.GetUserPosts(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "获取用户帖子失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": data,
	})
}
