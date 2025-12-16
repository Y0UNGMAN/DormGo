package controller

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/Y0UNGMAN/DormGo/backend/logic"
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
	data, err := logic.GetPostDetail(id)
	if err != nil {
		fmt.Println(err)
		c.JSON(400, gin.H{
			"code":    "400",
			"message": "get post detail failed",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "get post detail success",
		"data":    data,
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
	data, err := logic.GetPosts()
	if err != nil {
		fmt.Println(err)
		c.JSON(400, gin.H{
			"code":    "400",
			"message": "get posts failed",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "get posts success",
		"data":    data,
	})
}
