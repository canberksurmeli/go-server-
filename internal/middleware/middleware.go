package middleware

import (
	"message-provider-go/internal/database"

	"github.com/gin-gonic/gin"
)

func GinCORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func DBMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		db := database.Get()
		c.Set("dbPool", db.Pool)
		c.Next()
	}
}
