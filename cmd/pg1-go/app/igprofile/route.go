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
	router.GET("/igprofiles/count", func(c *gin.Context) {
		countIgProfileHandler(c, StatusActive)
	})
	router.GET("/igprofile/:ig_id", getIgProfileHandler)

	reqAdmin := router.Group("")
	reqAdmin.Use(auth.RequiredAdmin())
	{
		reqAdmin.GET("/igprofiles/search_all", findIgProfileAllStatusHandler)
		reqAdmin.GET("/igprofiles/count_all", func(c *gin.Context) {
			countIgProfileHandler(c, StatusAll)
		})
		reqAdmin.POST("/igprofiles", newIgProfileHandler)
		reqAdmin.PATCH("/igprofile/:ig_id", modifyIgProfileHandler)
		reqAdmin.DELETE("/igprofile/:ig_id", func(c *gin.Context) {
			deleteIgProfileHandler(c, false)
		})

		reqAdmin.POST("/igprofiles/action", igProfileActionHandler)

		reqAdmin.GET("/multi_acc_count", func(c *gin.Context) {
			countIgProfileHandler(c, StatusMulti)
		})
		reqAdmin.GET("/multi_acc", findMultiAccHandler)
		reqAdmin.POST("/multi_acc/:ig_id", activateMultiAccHandler)
		reqAdmin.DELETE("/multi_acc/:ig_id", func(c *gin.Context) {
			deleteIgProfileHandler(c, true)
		})
	}
}

// DefineViewRoutes defines routes for IgProfile that contains view
func DefineViewRoutes(router *gin.Engine, prefix string) {
	router.GET(prefix+"/", func(c *gin.Context) {
		igProfilesView(c, false)
	})
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
		adminPage.GET("/igprofiles", func(c *gin.Context) {
			igProfilesView(c, true)
		})
	}
}
