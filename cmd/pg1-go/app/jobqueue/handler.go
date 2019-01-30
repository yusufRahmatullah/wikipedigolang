package jobqueue

import (
	"net/http"

	"git.heroku.com/pg1-go-work/cmd/pg1-go/app/base"

	"github.com/gin-gonic/gin"
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
		if suc {
			msg = "Create JobQueue successful"
			status = http.StatusCreated
		} else {
			msg = "Failed to create JobQueue"
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
			for _, igID := range ids.([]interface{}) {
				ijq := NewJobQueue("SingleAccountJob", gin.H{"ig_id": igID.(string)})
				Save(ijq)
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
