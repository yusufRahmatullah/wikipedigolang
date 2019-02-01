package job

import (
	"git.heroku.com/pg1-go-work/cmd/pg1-go/app/igprofile"
	"git.heroku.com/pg1-go-work/cmd/pg1-go/app/jobqueue"
)

// BanAccountJob is the job to ban an account
type BanAccountJob struct{}

// NewBanAccountJob instantiate BanAccountJob instance
func NewBanAccountJob() *BanAccountJob {
	return &BanAccountJob{}
}

// Name reutrn BanAccountJob name
func (job *BanAccountJob) Name() string {
	return "BanAccountJob"
}

// Process executes job queue with the given params
// Returns true if process success
func (job *BanAccountJob) Process(jq *jobqueue.JobQueue) bool {
	params := (*jq).Params
	igID, ok := params["ig_id"]
	if ok {
		cleanID := igprofile.CleanIgIDParams(igID.(string))
		return igprofile.DeleteIgProfile(cleanID)
	}
	return false
}
