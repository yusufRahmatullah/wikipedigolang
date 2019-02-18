package igprofile

import (
	"fmt"
	"net/http"
	"strconv"

	"git.heroku.com/pg1-go-work/cmd/pg1-go/app/utils"

	"git.heroku.com/pg1-go-work/cmd/pg1-go/app/jobqueue"

	"github.com/globalsign/mgo/bson"

	"git.heroku.com/pg1-go-work/cmd/pg1-go/app/logger"

	"git.heroku.com/pg1-go-work/cmd/pg1-go/app/base"

	"github.com/gin-gonic/gin"
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
		if suc == "" {
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
	fr := utils.GetFindRequest(c)
	igps := GetAll(fr, StatusActive)
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

func modifyIgProfileHandler(c *gin.Context) {
	var igp IgProfile
	c.BindJSON(&igp)
	igID := c.Param("ig_id")
	changes := GenerateChanges(&igp)

	suc := Update(igID, changes)
	var msg string
	var status int
	if suc == "" {
		msg = "Update IgProfile successful"
		status = http.StatusCreated
	} else {
		msg = "Failed to update IgProfile"
		status = http.StatusOK
	}
	data := base.StandardJSON(msg, nil)
	c.JSON(status, data)
}

func convertIntOrDefault(text string, def int) int {
	if text == "" {
		return def
	}
	num, err := strconv.Atoi(text)
	if err != nil {
		handlerLogger.Warning(fmt.Sprintf("Failed to convert text to int: %v", text))
		return def
	}
	return num
}

func findIgProfile(c *gin.Context, status ProfileStatus) {
	fr := utils.GetFindRequest(c)
	igps := FindIgProfile(fr, status)
	compData := struct {
		Profiles []IgProfile `json:"profiles"`
		Query    string      `json:"query"`
	}{
		Profiles: igps,
		Query:    fr.Query,
	}
	data := base.StandardJSON("", compData)
	c.JSON(http.StatusOK, data)
}

func findIgProfileAllStatusHandler(c *gin.Context) {
	filterStatus := c.Query("filterStatus")
	status := StatusAll
	switch filterStatus {
	case "active":
		status = StatusActive
	case "banned":
		status = StatusBanned
	case "multi":
		status = StatusMulti
	}
	findIgProfile(c, status)
}

func findIgProfileHandler(c *gin.Context) {
	findIgProfile(c, StatusActive)
}

func deleteIgProfileHandler(c *gin.Context, multi bool) {
	igID := c.Param("ig_id")
	suc := DeleteIgProfile(igID, multi)
	var msg string
	if suc == "" {
		msg = "Delete IgProfile successful"
	} else {
		msg = "Failed to delete IgProfile"
	}
	data := base.StandardJSON(msg, nil)
	c.JSON(http.StatusOK, data)
}

func activateMultiAccHandler(c *gin.Context) {
	iid := c.Param("ig_id")
	suc := Update(iid, bson.M{"status": StatusMulti})
	var msg string
	if suc == "" {
		msg = "Save MultiAcc successful"
	} else {
		msg = "Failed to save MultiAcc"
	}
	data := base.StandardJSON(msg, nil)
	c.JSON(http.StatusOK, data)
}

func findMultiAccHandler(c *gin.Context) {
	fr := utils.GetFindRequest(c)
	mas := FindIgProfile(fr, StatusMulti)
	compData := struct {
		Accounts []IgProfile `json:"accounts"`
		Query    string      `json:"query"`
	}{
		Accounts: mas,
		Query:    fr.Query,
	}
	data := base.StandardJSON("", compData)
	c.JSON(http.StatusOK, data)
}

func getStatusFromAction(action string) string {
	status := "active"
	switch action {
	case "ban":
		status = "banned"
	case "asMulti":
		status = "multi"
	case "activate":
		status = "active"
	default:
		status = "active"
	}
	return status
}

func igProfileActionHandler(c *gin.Context) {
	jd := struct {
		IGID   string `json:"ig_id"`
		Action string `json:"action"`
	}{
		IGID:   "",
		Action: "",
	}
	c.BindJSON(&jd)
	if jd.IGID == "" {
		data := base.ErrorJSON("Param ig_id can't be empty", nil)
		c.JSON(http.StatusBadRequest, data)
	}
	status := getStatusFromAction(jd.Action)
	suc := Update(jd.IGID, bson.M{"status": status})
	if suc == "" {
		data := base.StandardJSON(fmt.Sprintf("Success to %v IG ID: %v", jd.Action, jd.IGID), nil)
		jq := jobqueue.NewJobQueue("UpdateIgMediaStatusJob", map[string]interface{}{"ig_id": jd.IGID})
		jobqueue.Save(jq)
		c.JSON(http.StatusOK, data)
	} else {
		data := base.ErrorJSON(fmt.Sprintf("Failed to %v IG ID: %v", jd.Action, jd.IGID), nil)
		c.JSON(http.StatusBadRequest, data)
	}
}

func countIgProfileHandler(c *gin.Context, status ProfileStatus) {
	n, err := countIgProfiles(status)
	if err == nil {
		data := base.StandardJSON("Success to count all IgProfiles", n)
		c.JSON(http.StatusOK, data)
	} else {
		data := base.ErrorJSON("Failed to count all IgProfiles", err.Error())
		c.JSON(http.StatusInternalServerError, data)
	}
}

/////////////////////////////////
// IgProfile Views
/////////////////////////////////

func igProfilesView(c *gin.Context, admin bool) {
	c.HTML(http.StatusOK, "igpage.tmpl.html", gin.H{
		"admin":     admin,
		"title":     "IGO Profile",
		"search_ph": "Search IGO here...",
		"sort_criteria": [][]string{
			[]string{"created_at", "Created At"},
			[]string{"modified_at", "Modified At"},
			[]string{"ig_id", "IG ID"},
			[]string{"name", "Name"},
			[]string{"followers", "Followers"},
			[]string{"following", "Following"},
			[]string{"posts", "Post Number"},
		},
		"filter_criteria": [][]string{
			[]string{"", "All"},
			[]string{"active", "Active"},
			[]string{"banned", "Banned"},
			[]string{"multi", "Multi"},
		},
		"script": gin.H{
			"name":         "igprofile",
			"page":         "igprofiles",
			"subtitle":     "Profiles",
			"search":       "/search",
			"count":        "/count",
			"admin_search": "/search_all",
			"admin_count":  "/count_all",
		},
	})
}

func multiAccView(c *gin.Context) {
	c.HTML(http.StatusOK, "multi_acc.tmpl.html", nil)
}
