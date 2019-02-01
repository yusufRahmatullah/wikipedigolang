package igprofile

import (
	"git.heroku.com/pg1-go-work/cmd/pg1-go/app/auth"
	"github.com/gin-gonic/gin"
)

// DefineAPIRoutes defines routes for Ig Profile
func DefineAPIRoutes(router *gin.Engine, prefix string) {
	router.GET(prefix+"/igprofiles", getAllIgProfileHandler)
	router.GET(prefix+"/igprofiles/search", findIgProfileHandler)
	router.GET(prefix+"/igprofile/:ig_id", getIgProfileHandler)

	reqAdmin := router.Group(prefix)
	reqAdmin.Use(auth.RequiredAdmin())
	{
		reqAdmin.POST(prefix+"/igprofiles", newIgProfileHandler)
		reqAdmin.PATCH(prefix+"/igprofile/:ig_id", modifyIgProfileHandler)
		reqAdmin.DELETE(prefix+"/igprofile/:ig_id", deleteIgProfileHandler)

		reqAdmin.GET("/multi_acc", findMultiAccHandler)
		reqAdmin.POST("/multi_acc/:ig_id", activateMultiAccHandler)
		reqAdmin.DELETE("/multi_acc/:ig_id", deleteMultiAccHandler)
	}
}

// DefineViewRoutes defines routes for IgProfile that contains view
func DefineViewRoutes(router *gin.Engine, prefix string) {
	router.GET(prefix+"/igprofiles", igProfilesView)

	reqAdmin := router.Group(prefix)
	reqAdmin.Use(auth.RequiredAdmin())
	{
		reqAdmin.GET("/multi_acc", multiAccView)
	}
}
