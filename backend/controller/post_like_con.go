package controller

import (
	"net/http"

	"github.com/Y0UNGMAN/DormGo/backend/logic"
	"github.com/Y0UNGMAN/DormGo/backend/model"
	"github.com/gin-gonic/gin"
)

type PostLikeRequest struct {
	PostID  uint `json:"postid" binding:"required"`
	LikerID uint `json:"likerid" binding:"required"`
}

func PostLike(c *gin.Context) {
	l := new(PostLikeRequest)
	err := c.ShouldBindJSON(l)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
			"msg":  "点赞信息解析失败",
		})
		return
	}
	likeDetail := &model.DgPostLike{
		PostID:  l.PostID,
		LikerID: l.LikerID,
	}
	err = logic.PostLike(likeDetail)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": http.StatusInternalServerError,
			"msg":  "PostLike error",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"msg":  "PostLike success",
	})
}

func CancelPostLike(c *gin.Context) {
	l := new(PostLikeRequest)
	err := c.ShouldBindJSON(l)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
			"msg":  "参数解析失败",
		})
		return
	}

	err = logic.CancelPostLike(l.PostID, l.LikerID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": http.StatusInternalServerError,
			"msg":  "CancelPostLike error",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"msg":  "CancelPostLike success",
	})

}
