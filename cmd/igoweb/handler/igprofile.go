package handler

import (
	"net/http"

	"git.heroku.com/pg1-go-work/cmd/igoweb/service"

	"git.heroku.com/pg1-go-work/cmd/igoweb/model"
	"git.heroku.com/pg1-go-work/cmd/igoweb/repository"
	"github.com/gin-gonic/gin"
)

// IgProfileHandler handles request about IgProfile
type IgProfileHandler struct {
	Repository repository.IgProfileRepository
	Logger     *service.LoggerService
}

// CreateIgProfile handles request to create new IgProfile
func (handler *IgProfileHandler) CreateIgProfile(c *gin.Context) {
	var igProfile *model.IgProfile
	c.BindJSON(igProfile)
	if igProfile.IGID == "" {
		resp := model.ErrorJSON("Param ig_id is required", nil)
		c.JSON(http.StatusBadRequest, resp)
	} else {
		err := handler.Repository.Create(igProfile)
		if err != nil {
			handler.Logger.Fatal("Failed to create IgProfile", err)
			resp := model.ErrorJSON("Failed to create IgProfile", err)
			c.JSON(http.StatusBadRequest, resp)
		} else {
			resp := model.SuccessJSON("Create IgProfile successful", nil)
			c.JSON(http.StatusCreated, resp)
		}
	}
}

// UpdateIgProfile handles request to update IgProfile
func (handler *IgProfileHandler) UpdateIgProfile(c *gin.Context) {
	var igProfile *model.IgProfile
	c.BindJSON(igProfile)
	if igProfile.IGID == "" {
		resp := model.ErrorJSON("Param ig_id is required", nil)
		c.JSON(http.StatusBadRequest, resp)
	} else {
		err := handler.Repository.Update(igProfile)
		if err != nil {
			handler.Logger.Fatal("Failed to update IgProfile", err)
			resp := model.ErrorJSON("Failed to update IgProfile", err)
			c.JSON(http.StatusBadRequest, resp)
		} else {
			resp := model.SuccessJSON("Update IgProfile successful", nil)
			c.JSON(http.StatusCreated, resp)
		}
	}
}

// FindIgProfiles handles request to find IgProfiles with given criteria
func (handler *IgProfileHandler) FindIgProfiles(c *gin.Context, status model.ProfileStatus) {
	fr := model.GetFindRequest(c)
	igProfiles, err := handler.Repository.Find(fr, status)
	if err != nil {
		handler.Logger.Fatal("Failed to find igprofiles", err)
		resp := model.ErrorJSON("Failed to find igprofiles", err)
		c.JSON(http.StatusBadRequest, resp)
	} else {
		resp := model.SuccessJSON("", igProfiles)
		c.JSON(http.StatusOK, resp)
	}
}

// CountIgProfiles handles request to count IgProfiles with given criteria
func (handler *IgProfileHandler) CountIgProfiles(c *gin.Context, status model.ProfileStatus) {
	fr := model.GetFindRequest(c)
	count, err := handler.Repository.Count(fr, status)
	if err != nil {
		handler.Logger.Fatal("Failed to count igprofiles", err)
		resp := model.ErrorJSON("Failed to count igprofiles", err)
		c.JSON(http.StatusBadRequest, resp)
	} else {
		resp := model.SuccessJSON("", count)
		c.JSON(http.StatusOK, resp)
	}
}
