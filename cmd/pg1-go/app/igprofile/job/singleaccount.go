package job

import (
	"strings"

	"git.heroku.com/pg1-go-work/cmd/pg1-go/app/igprofile"
	"git.heroku.com/pg1-go-work/cmd/pg1-go/app/jobqueue"
	"git.heroku.com/pg1-go-work/cmd/pg1-go/app/logger"
)

var (
	commaVanisher = strings.NewReplacer(",", "")
	sajLogger     = logger.NewLogger("SingleAccountJob", true, true)
)

// SingleAccountJob is the job for crawling an account
// of IgProfile
type SingleAccountJob struct{}

// NewSingleAccountJob instantiate SingleAccountJob instance
func NewSingleAccountJob() *SingleAccountJob {
	return &SingleAccountJob{}
}

// Name returns the SingleAccountJob
func (job *SingleAccountJob) Name() string {
	return "SingleAccountJob"
}

func crawlIgID(igID string) string {
	igp := igprofile.FetchIgProfile(igID)
	success := ""
	if igp != nil && igp.Posts > 0 {
		success = igprofile.SaveOrUpdate(igp)
		if success == "" {
			igp = igprofile.GetIgProfile(igID)
			if igp.Status == igprofile.StatusActive {
				jq := jobqueue.NewJobQueue("PostMediaJob", map[string]interface{}{"ig_id": igID})
				success = jobqueue.Save(jq)
			}
		}
	}
	return success
}

// Process executes job queue with the given params
// Returns empty string if process success
func (job *SingleAccountJob) Process(jq *jobqueue.JobQueue) string {
	sajLogger.Debug("run process")
	params := (*jq).Params
	igID, ok := params["ig_id"]
	if ok {
		cleanID := igprofile.CleanIgIDParams(igID.(string))
		return crawlIgID(cleanID)
	}
	return "Param ig_id not found"
}
