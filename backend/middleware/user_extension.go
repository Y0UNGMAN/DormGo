package middleware

import "github.com/gin-gonic/gin"

// Placeholder middleware file to keep package consistent.
func AllowAllCors() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
	}
}
