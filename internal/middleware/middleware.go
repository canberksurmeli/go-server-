package middleware

import (
	"message-provider-go/internal/database"

	"github.com/gin-gonic/gin"
)

// GinLogging returns a gin.HandlerFunc for logging requests
// func GinLogging(logger *logrus.Logger) gin.HandlerFunc {
// 	return gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
// 		logger.WithFields(logrus.Fields{
// 			"client_ip":   param.ClientIP,
// 			"method":      param.Method,
// 			"path":        param.Path,
// 			"status_code": param.StatusCode,
// 			"latency":     param.Latency,
// 			"user_agent":  param.Request.UserAgent(),
// 		}).Info("Request processed")
// 		return ""
// 	})
// }

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
