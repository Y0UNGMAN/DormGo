package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// 解决前后端分离项目中的跨域问题（CORS: Cross-Origin Resource Sharing）
func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		c.Header("Access-control-Allow-Origin", "*")
		c.Header("Access-control-Allow-Methods", "POST,GET,OPTIONS,PUT,DELETE,UPDATE")
		c.Header("Access-control-Allow-Headers", "origin,x-Requested-with,Content-Type,Accept,AccessToken,userToken,Authorization")
		c.Header("Access-control-Expose-Headers", "Content-Length,Access-Control-Allow-Origin,Access-Control-Allow-Headers,Cache-Control,Content-Language,Content-Type")
		c.Header("Access-Control-Allow-Credentials", "true")

		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
		}

		c.Next()
	}
}
