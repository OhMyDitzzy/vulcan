package api

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

// RequestLogger is a middleware that logs HTTP requests.
// We log the method, path, status code, and duration for observability.
func RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		method := c.Request.Method

		c.Next()

		duration := time.Since(start)
		statusCode := c.Writer.Status()
		
		log.Printf("[%s] %s - Status: %d - Duration: %v",
			method,
			path,
			statusCode,
			duration,
		)
	}
}