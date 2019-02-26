package jobqueue

import (
	"fmt"
	"net/http"

	"git.heroku.com/pg1-go-work/cmd/pg1-go/app/utils"

	"git.heroku.com/pg1-go-work/cmd/pg1-go/app/base"

	"github.com/gin-gonic/gin"
	"github.com/globalsign/mgo/bson"
)

func getAvailableJobsHandler(c *gin.Context) {
	avaJobs := getAvailableJobs()
	data := base.StandardJSON("", avaJobs)
	c.JSON(http.StatusOK, data)
}

// newJobQueueHandler handles the request to post a new JobQueue
func newJobQueueHandler(c *gin.Context) {
	var jq JobQueue
	var data interface{}
	c.BindJSON(&jq)
	if jq.Name == "" {
		data = base.ErrorJSON("name is required", nil)
		c.JSON(http.StatusBadRequest, data)
	} else {
		suc := Save(&jq)
		var msg string
		var status int
		if suc == "" {
			msg = "Create JobQueue successful"
			status = http.StatusCreated
		} else {
			msg = suc
			status = http.StatusOK
		}
		data := base.StandardJSON(msg, nil)
		c.JSON(status, data)
	}
}

func batchAddHandler(c *gin.Context) {
	var jq JobQueue
	c.BindJSON(&jq)
	if jq.Name == "" {
		data := base.ErrorJSON("name is required", nil)
		c.JSON(http.StatusBadRequest, data)
	} else {
		ids, ok := jq.Params["ig_ids"]
		if ok {
			jobName, ok := jq.Params["job_name"]
			if ok {
				for _, igID := range ids.([]interface{}) {
					ijq := NewJobQueue(jobName.(string), gin.H{"ig_id": igID.(string)})
					Save(ijq)
				}
			}
		}
		data := base.StandardJSON("Batch Add JobQueue successful", nil)
		c.JSON(http.StatusOK, data)
	}
}

func getPostponedJobsHandler(c *gin.Context) {
	jqs := GetAllPostponed()
	data := base.StandardJSON("", jqs)
	c.JSON(http.StatusOK, data)
}

func deletePostponedHandler(c *gin.Context) {
	jqid := c.Param("job_id")
	jq := GetPostponed(jqid)
	suc := DeletePostponed(jq)
	if suc {
		data := base.StandardJSON("Delete postponed JobQueue success", nil)
		c.JSON(http.StatusOK, data)
	} else {
		data := base.ErrorJSON("Failed to delete postponed JobQueue", nil)
		c.JSON(http.StatusNotModified, data)
	}
}

func requeuePostponedHandler(c *gin.Context) {
	jdata := struct {
		JobID string `json:"job_id"`
	}{
		JobID: "",
	}
	c.BindJSON(&jdata)
	if jdata.JobID == "" {
		data := base.ErrorJSON("Job ID is empty", nil)
		c.JSON(http.StatusNotModified, data)
	} else {
		jq := GetPostponed(jdata.JobID)
		if jq.ID != "" {
			suc := RequeuePostponed(jq)
			if suc {
				data := base.StandardJSON("Requeue postponed JobQueue success", nil)
				c.JSON(http.StatusOK, data)
			} else {
				data := base.ErrorJSON("Failed to requeue postponed JobQueue", nil)
				c.JSON(http.StatusNotModified, data)
			}
		} else {
			data := base.ErrorJSON("Failed to get postponed JobQueue", nil)
			c.JSON(http.StatusNotModified, data)
		}
	}
}

func countPostponedJobsHandler(c *gin.Context) {
	n, err := countPostponedJobs()
	if err == nil {
		data := base.StandardJSON("", n)
		c.JSON(http.StatusOK, data)
	} else {
		data := base.ErrorJSON("", err.Error())
		c.JSON(http.StatusInternalServerError, data)
	}
}

func findJobQueueHandler(c *gin.Context) {
	filterStatus := c.Query("filterStatus")
	status := StatusActive
	switch filterStatus {
	case "postponed":
		status = StatusPostponed
	case "finished":
		status = StatusFinished
	case "":
		status = StatusAll
	}
	fr := utils.GetFindRequest(c)
	jqs := FindJobQueue(fr, status)
	compData := struct {
		Jobs  []JobQueue `json:"jobs"`
		Query string     `json:"query"`
	}{
		Jobs:  jqs,
		Query: fr.Query,
	}
	data := base.StandardJSON("", compData)
	c.JSON(http.StatusOK, data)
}

func actionHandler(c *gin.Context) {
	jd := struct {
		JobID  string `json:"job_id"`
		Action string `json:"action"`
	}{
		JobID:  "",
		Action: "",
	}
	c.BindJSON(&jd)
	if jd.JobID == "" {
		data := base.ErrorJSON("Param job_id can't be empty", nil)
		c.JSON(http.StatusBadRequest, data)
		return
	}
	status := StatusActive
	if jd.Action == "delete" {
		status = StatusFinished
	}
	jq := GetJobQueue(jd.JobID)
	suc := Update(jq, bson.M{"status": status})
	if suc {
		data := base.StandardJSON(fmt.Sprintf("Success to %v Job ID: %v", jd.Action, jd.JobID), nil)
		c.JSON(http.StatusBadRequest, data)
	} else {
		data := base.ErrorJSON(fmt.Sprintf("Failed to %v Job ID: %v", jd.Action, jd.JobID), nil)
		c.JSON(http.StatusNotModified, data)
	}
}

// jobQueueIndexView render JobQueue form
func jobQueueIndexView(c *gin.Context) {
	c.HTML(http.StatusOK, "jobqueue.tmpl.html", nil)
}

func batchAddIndexView(c *gin.Context) {
	c.HTML(http.StatusOK, "batch_add.tmpl.html", nil)
}

func postponedJobsView(c *gin.Context) {
	c.HTML(http.StatusOK, "postponed_jobs.tmpl.html", nil)
}
