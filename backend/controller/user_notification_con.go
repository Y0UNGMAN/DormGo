package controller

import (
	"fmt"

	"github.com/Y0UNGMAN/DormGo/backend/logic"
	"github.com/Y0UNGMAN/DormGo/backend/middleware"
	"github.com/Y0UNGMAN/DormGo/backend/model"
	"github.com/gin-gonic/gin"
)

func GetNotifacation(c *gin.Context) {
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
	var notifs []model.DgNotification

	notifs, err := logic.GetNotifacation(userId)

	if err != nil {
		fmt.Println(err)
		c.JSON(500, gin.H{"code": 500, "msg": "获取通知失败"})
		return
	}

	c.JSON(200, gin.H{
		"code": 200,
		"data": notifs,
	})

}

func GetUnreadCount(c *gin.Context) {
	var userId uint = 0
	if v, exists := c.Get(middleware.CtxUserIDKey); exists {
		if uid, ok := v.(uint); ok {
			userId = uid
		} else if uidInt, ok := v.(int); ok {
			userId = uint(uidInt)
		}
	}
	count, err := logic.GetUnreadCount(userId)
	if err != nil {
		c.JSON(500, gin.H{"code": 500, "msg": "获取未读数量失败"})
		return
	}
	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "Success",
		"data": count,
	})
}

func ReadAllNotification(c *gin.Context) {
	var userId uint = 0
	if v, exists := c.Get(middleware.CtxUserIDKey); exists {
		if uid, ok := v.(uint); ok {
			userId = uid
		} else if uidInt, ok := v.(int); ok {
			userId = uint(uidInt)
		}
	}
	err := logic.ReadAllNotification(userId)
	if err != nil {
		fmt.Println(err)
		c.JSON(500, gin.H{
			"code": 500,
			"msg":  "failed",
		})
		return
	}
	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "已读所有通知",
	})
}
