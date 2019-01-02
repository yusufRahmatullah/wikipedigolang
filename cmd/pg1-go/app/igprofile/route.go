package igprofile

import (
	"github.com/gin-gonic/gin"
)

// DefineAPIRoutes defines routes for Ig Profile
func DefineAPIRoutes(router *gin.Engine, prefix string) {
	router.GET(prefix+"/igprofiles", getAllIgProfileHandler)
	router.POST(prefix+"/igprofiles", newIgProfileHandler)
	router.GET(prefix+"/igprofiles/search", findIgProfileHandler)
	router.GET(prefix+"/igprofile/:ig_id", getIgProfileHandler)
	router.PATCH(prefix+"/igprofile/:ig_id", modifyIgProfileHandler)
	router.DELETE(prefix+"/igprofile/:ig_id", deleteIgProfileHandler)
}

// DefineViewRoutes defines routes for IgProfile that contains view
func DefineViewRoutes(router *gin.Engine, prefix string) {
	router.GET(prefix+"/igprofiles", igProfilesView)
}
