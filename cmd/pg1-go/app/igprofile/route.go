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
		reqAdmin.GET("/igprofiles/search_all", findIgProfileAllStatusHandler)
		reqAdmin.POST("/igprofiles", newIgProfileHandler)
		reqAdmin.PATCH("/igprofile/:ig_id", modifyIgProfileHandler)
		reqAdmin.DELETE("/igprofile/:ig_id", deleteIgProfileHandler)

		reqAdmin.POST("/igprofiles/action", igProfileActionHandler)

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

	adminPage := router.Group("/admin")
	adminPage.Use(auth.RequiredAdmin())
	{
		adminPage.GET("/igprofiles", adminIgProfileView)
	}
}
