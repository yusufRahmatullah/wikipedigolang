package repeat

import (
	"github.com/gin-gonic/gin"
)

// DefineRoutes defines repeat-specific routes
func DefineRoutes(router *gin.Engine) {
	router.GET("/repeat", Handler)
}
