package middleware

import (
	"github.com/gin-gonic/gin"
)

// NoCacheHeader is a gin middleware that used to set cache control to headers
func NoCacheHeader() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
		c.Writer.Header().Set("Pragma", "no-cache")
		c.Writer.Header().Set("Expires", "0")
		c.Next()
	}
}
