package job

import (
	"fmt"

	"git.heroku.com/pg1-go-work/cmd/pg1-go/app/utils"

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

func updateIgID(igp *igprofile.IgProfile) string {
	igID := igp.IGID
	igp2 := igprofile.FetchIgProfile(igID)
	if igp2 != nil {
		changes := igprofile.GenerateChanges(igp2)
		suc := igprofile.Update(igID, changes)
		if suc == "" {
			ujLogger.Debug(fmt.Sprintf("Success to update IG ID: %v", igID))
		} else {
			ujLogger.Fatal(fmt.Sprintf("Failed to update IG ID: %v", igID), nil)
			return suc
		}
		if igp2.Status == igprofile.StatusActive {
			jq := jobqueue.NewJobQueue("PostMediaJob", map[string]interface{}{
				"ig_id": igID,
			})
			return jobqueue.Save(jq)
		} else if igp2.Status == igprofile.StatusMulti {
			jq := jobqueue.NewJobQueue("PostAccountJob", map[string]interface{}{
				"ig_id": igID,
			})
			return jobqueue.Save(jq)
		}
	}
	return fmt.Sprintf("Failed to fetch IgProfile with IG ID: %v", igp.IGID)
}

// Process executes job queue with the given params
// Update process is not guaranteed to success
// This method always returns empty string
func (job *UpdaterJob) Process(jq *jobqueue.JobQueue) string {
	ujLogger.Debug("run process")
	ujLogger.Debug("Get all existing IgProfile")
	var igps []igprofile.IgProfile
	offset := 0
	limit := 10
	fr := utils.FindRequest{
		Offset: 0, Limit: 10,
	}
	igps = igprofile.GetAll(&fr, igprofile.StatusAll)
	for len(igps) > 0 {
		ujLogger.Debug(fmt.Sprintf("offset: %v, limit: %v, len(igps): %v", offset, limit, len(igps)))
		for i, igp := range igps {
			ujLogger.Debug(fmt.Sprintf("Processing item: %v", (i + offset)))
			updateIgID(&igp)
		}
		// update igps and offset
		if len(igps) == limit {
			fr.Offset += limit
			igps = igprofile.GetAll(&fr, igprofile.StatusActive)
		} else {
			break
		}
	}
	return ""
}

// SingleUpdaterJob is the job for update single IgProfile
type SingleUpdaterJob struct{}

// NewSingleUpdaterJob instantiate SingleUpdaterJob instance
func NewSingleUpdaterJob() *SingleUpdaterJob {
	return &SingleUpdaterJob{}
}

// Name returns "SingleUpdaterJob"
func (job *SingleUpdaterJob) Name() string {
	return "SingleUpdaterJob"
}

// Process executes job queue with the given params
// Update process is not guaranteed to success
// Returns empty string if success otherwise
// returns error string
func (job *SingleUpdaterJob) Process(jq *jobqueue.JobQueue) string {
	igID, ok := jq.Params["ig_id"]
	if ok {
		cleanID := igprofile.CleanIgIDParams(igID.(string))
		igp := igprofile.GetIgProfile(igID.(string))
		if igp == nil {
			return fmt.Sprintf("Failed to get IgProfile with IG ID: %v", cleanID)
		}
		return updateIgID(igp)
	}
	return "Param ig_id not found"
}
