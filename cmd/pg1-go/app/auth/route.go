package auth

import "github.com/gin-gonic/gin"

// DefineViewRoutes defines routes for authentication
func DefineViewRoutes(router *gin.Engine, prefix string) {
	router.GET(prefix+"/login", loginView)
	router.GET(prefix+"/logout", logout)
}

// DefineAPIRoutes defines routes for authentication API
func DefineAPIRoutes(router *gin.RouterGroup) {
	router.POST("/login", login)
}
