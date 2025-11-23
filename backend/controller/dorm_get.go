package controller

import (
	"fmt"
	"net/http"

	"github.com/Y0UNGMAN/DormGo/backend/logic"
	"github.com/gin-gonic/gin"
)

func GetDorms(c *gin.Context) {
	data, err := logic.GetDorms()
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
		})
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "success",
		"data": data,
	})
}
