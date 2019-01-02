package igprofile

import (
	"fmt"
	"net/http"
	"strconv"

	"git.heroku.com/pg1-go-work/cmd/pg1-go/app/logger"

	"git.heroku.com/pg1-go-work/cmd/pg1-go/app/base"

	"github.com/gin-gonic/gin"
)

const (
	defaultOffset = 0
	defaultLimit  = 20
)

var handlerLogger = logger.NewLogger("IgProfileHandler", false, true)

func newIgProfileHandler(c *gin.Context) {
	var igp IgProfile
	c.BindJSON(&igp)
	if igp.IGID == "" {
		data := base.ErrorJSON("ig_id is required", nil)
		c.JSON(http.StatusBadRequest, data)
	} else {
		suc := Save(&igp)
		var msg string
		var status int
		if suc {
			msg = "Create IgProfile successful"
			status = http.StatusCreated
		} else {
			msg = "Failed to create IgProfile"
			status = http.StatusOK
		}
		data := base.StandardJSON(msg, nil)
		c.JSON(status, data)
	}
}

func getAllIgProfileHandler(c *gin.Context) {
	offsetStr := c.Query("offset")
	limitStr := c.Query("limit")
	var offset, limit int
	var err error
	if offsetStr == "" {
		offset = defaultOffset
	} else {
		offset, err = strconv.Atoi(offsetStr)
		if err != nil {
			handlerLogger.Warning(fmt.Sprintf("Failed to convert offset text: %v", offsetStr))
			offset = defaultOffset
		}
	}
	if limitStr == "" {
		limit = defaultLimit
	} else {
		limit, err = strconv.Atoi(limitStr)
		if err != nil {
			handlerLogger.Warning(fmt.Sprintf("Failed to convert limit text: %v", limitStr))
			limit = defaultLimit
		}
	}
	igps := GetAll(offset, limit)
	data := base.StandardJSON("", igps)
	c.JSON(http.StatusOK, data)
}

func getIgProfileHandler(c *gin.Context) {
	igID := c.Param("ig_id")
	igp := GetIgProfile(igID)
	var data map[string]interface{}
	if igp.IGID == "" {
		data = base.StandardJSON("User not found", nil)
	} else {
		data = base.StandardJSON("", igp)
	}
	c.JSON(http.StatusOK, data)
}

func generateChanges(igp *IgProfile) *gin.H {
	changes := gin.H{}
	if igp.Name != "" {
		changes["name"] = igp.Name
	}
	if igp.Followers > 0 {
		changes["followers"] = igp.Followers
	}
	if igp.Following > 0 {
		changes["following"] = igp.Following
	}
	if igp.Posts > 0 {
		changes["posts"] = igp.Posts
	}
	return &changes
}

func modifyIgProfileHandler(c *gin.Context) {
	var igp IgProfile
	c.BindJSON(&igp)
	igID := c.Param("ig_id")
	changes := generateChanges(&igp)

	suc := Update(igID, *changes)
	var msg string
	var status int
	if suc {
		msg = "Update IgProfile successful"
		status = http.StatusCreated
	} else {
		msg = "Failed to update IgProfile"
		status = http.StatusOK
	}
	data := base.StandardJSON(msg, nil)
	c.JSON(status, data)
}

func findIgProfileHandler(c *gin.Context) {
	offsetStr := c.Query("offset")
	limitStr := c.Query("limit")
	query := c.Query("query")
	var offset, limit int
	var err error
	if offsetStr == "" {
		offset = defaultOffset
	} else {
		offset, err = strconv.Atoi(offsetStr)
		if err != nil {
			handlerLogger.Warning(fmt.Sprintf("Failed to convert offset text: %v", offsetStr))
			offset = defaultOffset
		}
	}
	if limitStr == "" {
		limit = defaultLimit
	} else {
		limit, err = strconv.Atoi(limitStr)
		if err != nil {
			handlerLogger.Warning(fmt.Sprintf("Failed to convert offset text: %v", limitStr))
			limit = defaultLimit
		}
	}
	igps := FindIgProfile(query, offset, limit)
	data := base.StandardJSON("", igps)
	c.JSON(http.StatusOK, data)
}

func deleteIgProfileHandler(c *gin.Context) {
	igID := c.Param("ig_id")
	suc := DeleteIgProfile(igID)
	var msg string
	if suc {
		msg = "Delete IgProfile successful"
	} else {
		msg = "Failed to delete User"
	}
	data := base.StandardJSON(msg, nil)
	c.JSON(http.StatusOK, data)
}
