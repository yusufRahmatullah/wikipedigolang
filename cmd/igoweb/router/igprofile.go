package router

import (
	"git.heroku.com/pg1-go-work/cmd/igoweb/handler"
	"git.heroku.com/pg1-go-work/cmd/igoweb/middleware"
	"git.heroku.com/pg1-go-work/cmd/igoweb/model"
	"git.heroku.com/pg1-go-work/cmd/igoweb/repository"
	"github.com/gin-gonic/gin"
)

// IgProfileRouter redirect user request to corresponding handler
type IgProfileRouter struct {
	Handler        *handler.IgProfileHandler
	UserRepository repository.UserRepository
}

// DefineAPIRoutes defines routes for API call that return JSON response
func (router *IgProfileRouter) DefineAPIRoutes(rg *gin.RouterGroup) {
	rg.GET("/igprofiles", func(c *gin.Context) {
		router.Handler.FindIgProfiles(c, model.StatusActive)
	})
	rg.GET("/igprofiles/count", func(c *gin.Context) {
		router.Handler.CountIgProfiles(c, model.StatusActive)
	})
	rg.GET("/igprofile/:ig_id", func(c *gin.Context) {
		panic("Not Implemented")
	})

	reqAdmin := rg.Group("/admin")
	reqAdmin.Use(middleware.RequireAdmin(router.UserRepository))
	{
		reqAdmin.GET("/igprofiles", func(c *gin.Context) {
			status := getStatusFromFilter(c)
			router.Handler.FindIgProfiles(c, status)
		})
		reqAdmin.GET("/igprofiles/count", func(c *gin.Context) {
			status := getStatusFromFilter(c)
			router.Handler.CountIgProfiles(c, status)
		})
	}
}

func getStatusFromFilter(c *gin.Context) model.ProfileStatus {
	filterStatus := c.Query("filterStatus")
	switch filterStatus {
	case "active":
		return model.StatusActive
	case "banned":
		return model.StatusBanned
	case "multi":
		return model.StatusMulti
	default:
		return model.StatusAll
	}
}
