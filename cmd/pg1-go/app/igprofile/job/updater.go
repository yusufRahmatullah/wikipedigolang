package job

import (
	"fmt"

	"git.heroku.com/pg1-go-work/cmd/pg1-go/app/igprofile"
	"git.heroku.com/pg1-go-work/cmd/pg1-go/app/jobqueue"
	"git.heroku.com/pg1-go-work/cmd/pg1-go/app/logger"
)

var (
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
	igp2 := igprofile.FetchIgProfile(igID)
	if igp2 != nil {
		changes := igprofile.GenerateChanges(igp2)
		suc := igprofile.Update(igID, changes)
		if suc == "" {
			ujLogger.Debug(fmt.Sprintf("Success to update IG ID: %v", igID))
		} else {
			ujLogger.Fatal(fmt.Sprintf("Failed to update IG ID: %v", igID), nil)
		}
		if igp2.Status == igprofile.StatusActive {
			jq := jobqueue.NewJobQueue("PostMediaJob", map[string]interface{}{
				"ig_id": igID,
			})
			jobqueue.Save(jq)
		} else if igp2.Status == igprofile.StatusMulti {
			jq := jobqueue.NewJobQueue("PostAccountJob", map[string]interface{}{
				"ig_id": igID,
			})
			jobqueue.Save(jq)
		}
	}
}

// Process executes job queue with the given params
// Update process is not guaranted to success
// This method always returns empty string
func (job *UpdaterJob) Process(jq *jobqueue.JobQueue) string {
	ujLogger.Debug("run process")
	ujLogger.Debug("Get all existing IgProfile")
	var igps []igprofile.IgProfile
	offset := 0
	limit := 10
	igps = igprofile.GetAll(offset, limit, igprofile.StatusAll)
	for len(igps) > 0 {
		ujLogger.Debug(fmt.Sprintf("offset: %v, limit: %v, len(igps): %v", offset, limit, len(igps)))
		for i, igp := range igps {
			ujLogger.Debug(fmt.Sprintf("Processing item: %v", (i + offset)))
			updateIgID(&igp)
		}
		// update igps and offset
		if len(igps) == limit {
			offset += limit
			igps = igprofile.GetAll(offset, limit, igprofile.StatusActive)
		} else {
			break
		}
	}
	return ""
}
