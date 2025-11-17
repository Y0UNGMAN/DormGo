package controller

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/Y0UNGMAN/DormGo/backend/logic"
	"github.com/gin-gonic/gin"
)

func PostType(c *gin.Context) {
	//查询到所有帖子类型，以列表形式返回
	data, err := logic.GetPostType()
	if err != nil {
		fmt.Println(err)
		//ResponseError(c, CodeServerBusy)
		return
	}
	//ResponseSuccess(c, data)
	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"msg":  "success",
		"data": data,
	})

}

func PostTypeDetial(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	data, err := logic.GetPostTypeById(id)
	if err != nil {
		fmt.Println(err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"msg":  "success",
		"data": data,
	})
}
