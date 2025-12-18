package controller

import (
	"net/http"
	"strconv"

	"github.com/Y0UNGMAN/DormGo/backend/logic"
	"github.com/gin-gonic/gin"
)

type CreateCommentReq struct {
	PostID      uint   `json:"post_id" binding:"required"`
	Content     string `json:"content" binding:"required"`
	CommenterID uint   `json:"commenter_id" binding:"required"`
	ParentID    uint   `json:"parent_id"`
}

func CreateComment(c *gin.Context) {
	var req CreateCommentReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "数据解析失败",
		})
		return
	}
	logicCommentParams := logic.PublishCommentParams{
		PostID:   req.PostID,
		UserID:   req.CommenterID,
		Content:  req.Content,
		ParentID: req.ParentID,
	}

	err := logic.Comment.PublishComment(&logicCommentParams)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": http.StatusInternalServerError,
			"msg":  "logic 函数错误",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"msg":  "success",
	})

}

func GetComment(c *gin.Context) {
	idstr := c.Query("post_id")
	if idstr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "post_id 为必填项"})
		return
	}
	postID, _ := strconv.Atoi(idstr)

	pageStr := c.DefaultQuery("page", "1")
	page, _ := strconv.Atoi(pageStr)
	pageSize := 10

	resp, err := logic.Comment.GetPostComments(uint(postID), page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": http.StatusInternalServerError,
			"msg":  "获取评论失败",
		})
		return
	}

	hasMore := false
	if int64(page*pageSize) < resp.Total {
		hasMore = true
	}
	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"msg":  "success",
		"data": gin.H{
			"list":     resp.List,
			"total":    resp.Total,
			"page":     page,
			"has_more": hasMore,
		},
	})

}
