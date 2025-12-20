package controller

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/Y0UNGMAN/DormGo/backend/logic"
	"github.com/Y0UNGMAN/DormGo/backend/middleware"
	"github.com/Y0UNGMAN/DormGo/backend/model"
	"github.com/gin-gonic/gin"
)

func SendMessage(c *gin.Context) {
	var req model.SendMessageReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
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

	var msg model.DgMessage
	msg, err := logic.SendMessage(req, userId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 500,
			"msg":  "logic函数调用错误",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "success",
		"data": msg,
	})

}

func GetMessageHistory(c *gin.Context) {
	targetIDStr := c.Query("target_id")
	targetID, _ := strconv.Atoi(targetIDStr)
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
	go func() {
		err := logic.MarkAsRead(userId, uint(targetID))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code": 500,
				"msg":  "MarkAsRead",
			})
			return
		}
	}()

	var msgs []model.DgMessage
	msgs, err := logic.GetMessage(uint(targetID), userId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 500,
			"msg":  "logiC调用失败",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "success",
		"data": msgs,
	})

}

// 获取会话列表接口
func GetConversations(c *gin.Context) {
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
	if userId == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "msg": "未登录"})
		return
	}
	list, err := logic.GetConversationList(uint(userId))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "查询会话列表失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": list,
	})
}

// 获取总未读数接口
func GetUnreadMessageCount(c *gin.Context) {
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
	if userId == 0 {
		// 未登录返回0，不报错
		c.JSON(http.StatusOK, gin.H{"code": 200, "data": 0})
		return
	}

	count, err := logic.GetUnreadMessageCount(userId)
	if err != nil {
		c.JSON(500, gin.H{
			"code": 500,
			"msg":  "logic.GetUnreadCount错误",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "success",
		"data": count,
	})
}
