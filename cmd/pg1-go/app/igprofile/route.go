package igprofile

import (
	"net/http"

	"git.heroku.com/pg1-go-work/cmd/pg1-go/app/auth"
	"github.com/gin-gonic/gin"
)

// DefineAPIRoutes defines routes for Ig Profile
func DefineAPIRoutes(router *gin.RouterGroup) {
	router.GET("/igprofiles", getAllIgProfileHandler)
	router.GET("/igprofiles/search", findIgProfileHandler)
	router.GET("/igprofile/:ig_id", getIgProfileHandler)

	reqAdmin := router.Group("")
	reqAdmin.Use(auth.RequiredAdmin())
	{
		reqAdmin.POST("/igprofiles", newIgProfileHandler)
		reqAdmin.PATCH("/igprofile/:ig_id", modifyIgProfileHandler)
		reqAdmin.DELETE("/igprofile/:ig_id", deleteIgProfileHandler)

		reqAdmin.GET("/multi_acc", findMultiAccHandler)
		reqAdmin.POST("/multi_acc/:ig_id", activateMultiAccHandler)
		reqAdmin.DELETE("/multi_acc/:ig_id", deleteMultiAccHandler)
	}
}

// DefineViewRoutes defines routes for IgProfile that contains view
func DefineViewRoutes(router *gin.Engine, prefix string) {
	router.GET(prefix+"/", igProfilesView)
	router.GET(prefix+"/igprofiles", func(c *gin.Context) {
		c.Redirect(http.StatusPermanentRedirect, "/")
	})

	reqAdmin := router.Group(prefix)
	reqAdmin.Use(auth.RequiredAdmin())
	{
		reqAdmin.GET("/multi_acc", multiAccView)
	}
}
