package controller

import (
	"fmt"

	"github.com/Y0UNGMAN/DormGo/backend/logic"
	"github.com/Y0UNGMAN/DormGo/backend/model"
	"github.com/gin-gonic/gin"
)

func PostCreate(c *gin.Context) {
	//1 获取参数 参数校验
	p := new(model.DgPost)
	err := c.ShouldBindJSON(p)
	if err != nil {
		fmt.Println(err)
		c.JSON(400, gin.H{
			"code": 400,
			"msg":  "参数读取失败",
		})
		return
	}

	//2 创建帖子
	err = logic.CreatePost(p)
	if err != nil {
		fmt.Println(err)
		c.JSON(400, gin.H{
			"code": 400,
			"msg":  "创建帖子失败",
		})
		return
	}
	//3 返回参数

	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "Create Success",
	})

}
