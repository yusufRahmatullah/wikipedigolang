package job

import (
	"fmt"

	"github.com/globalsign/mgo/bson"

	"git.heroku.com/pg1-go-work/cmd/pg1-go/app/igprofile"
	"git.heroku.com/pg1-go-work/cmd/pg1-go/app/jobqueue"
	"git.heroku.com/pg1-go-work/cmd/pg1-go/app/logger"
	"git.heroku.com/pg1-go-work/cmd/pg1-go/app/utils"
)

var (
	updaterJsFunction = `function() {
		var ppNode = document.querySelector("span>img");
		if (ppNode) {
			var ppURL = ppNode.getAttribute("src");
		} else {
			ppNode = document.querySelector("button>img");
			var ppURL = "";
			if (ppNode) {
				ppURL = ppNode.getAttribute("src");
			}
		}
		return { ppURL: ppURL };
	}`
	ujLogger = logger.NewLogger("UpdaterJob", true, true)
)

// UpdaterJob is the job for update existing IgProfile
type UpdaterJob struct{}

// NewUpdaterJob instantiate UpdaterJob instance
func NewUpdaterJob() *UpdaterJob {
	return &UpdaterJob{}
}

// Name returns "UpdaterJob"
func (job *UpdaterJob) Name() string {
	return "UpdaterJob"
}

func updateIgID(igp *igprofile.IgProfile) {
	igID := igp.IGID
	wpw := utils.NewWebPageWrapper(ujLogger)
	if wpw != nil {
		defer wpw.Close()
		wpw.OnEvaluated(func(data map[string]interface{}) {
			ppURL := data["ppURL"].(string)
			ujLogger.Debug(fmt.Sprintf("ppURL: %v", ppURL))
			suc := igprofile.Update(igID, bson.M{"pp_url": ppURL})
			if suc {
				ujLogger.Debug(fmt.Sprintf("Success to update IG ID: %v", igID))
			} else {
				ujLogger.Fatal(fmt.Sprintf("Failed to update IG ID: %v", igID))
			}
		})
		wpw.OpenURL(fmt.Sprintf("https://www.instagram.com/%v", igID))
		wpw.Evaluate(updaterJsFunction)
	}
}

// Process executes job queue with the given params
// Update process is not guaranted to success
// This method always returns true
func (job *UpdaterJob) Process(jq *jobqueue.JobQueue) bool {
	ujLogger.Debug("run process")
	ujLogger.Debug("Get all existing IgProfile")
	var igps []igprofile.IgProfile
	offset := 0
	limit := 10
	igps = igprofile.GetAll(offset, limit)
	for len(igps) > 0 {
		ujLogger.Debug(fmt.Sprintf("offset: %v, limit: %v, len(igps): %v", offset, limit, len(igps)))
		for i, igp := range igps {
			ujLogger.Debug(fmt.Sprintf("Processing item: %v", (i + offset)))
			updateIgID(&igp)
		}
		// update igps and offset
		if len(igps) == limit {
			offset += limit
			igps = igprofile.GetAll(offset, limit)
		} else {
			break
		}
	}
	return true
}
