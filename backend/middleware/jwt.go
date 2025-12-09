package middleware

import (
	"strings"

	"github.com/Y0UNGMAN/DormGo/backend/utils"
	"github.com/gin-gonic/gin"
)

const CtxUserIDKey = "UserID"

func JWTAuthMiddleware() func(c *gin.Context) {
	return func(c *gin.Context) {
		authHeader := c.Request.Header.Get("Authorization")
		if authHeader == "" {
			c.JSON(200, gin.H{
				"code": 401,
				"msg":  "(未登录)请求头为空",
			})
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			c.JSON(200, gin.H{
				"code": 401,
				"msg":  "Invalid token",
			})
			c.Abort()
			return
		}

		mc, err := utils.ParseToken(parts[1])
		if err != nil {
			c.JSON(200, gin.H{
				"code": 401,
				"msg":  "无效token",
			})
			c.Abort()
			return
		}

		c.Set(CtxUserIDKey, mc.UserID)
		c.Next()

	}
}
