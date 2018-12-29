package user

import (
	"github.com/gin-gonic/gin"
)

// DefineRoutes defines User-specific routes
func DefineRoutes(router *gin.Engine, prefix string) {
	router.GET(prefix+"/users", usersHandler)
	router.POST(prefix+"/users", newUserHandler)
	router.GET(prefix+"/users/:name", findUserHandler)
	router.DELETE(prefix+"/users/:name", deleteUserHandler)
}
