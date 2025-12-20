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
	data, isliked, ispinned, isFavorited, err := logic.GetPostDetail(id, userId)
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
		"is_pinned":   ispinned,
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
