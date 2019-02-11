package job

import (
	"git.heroku.com/pg1-go-work/cmd/pg1-go/app/igmedia"
	"git.heroku.com/pg1-go-work/cmd/pg1-go/app/igprofile"
	"git.heroku.com/pg1-go-work/cmd/pg1-go/app/jobqueue"
	"git.heroku.com/pg1-go-work/cmd/pg1-go/app/logger"
)

var (
	usLogger = logger.NewLogger("UpdateIgMediaStatusJob", true, true)
)

// UpdateIgMediaStatusJob is the job for get top twelve post of IgProfile
type UpdateIgMediaStatusJob struct{}

// NewUpdateIgMediaStatusJob instantiate TopTwelveJob instance
func NewUpdateIgMediaStatusJob() *UpdateIgMediaStatusJob {
	return &UpdateIgMediaStatusJob{}
}

// Name returns UpdateIgMediaStatusJob
func (job *UpdateIgMediaStatusJob) Name() string {
	return "UpdateIgMediaStatusJob"
}

// Process executes JobQueue with the given params
// Returns empty string if success
func (job *UpdateIgMediaStatusJob) Process(jq *jobqueue.JobQueue) string {
	params := jq.Params
	igID, ok := params["ig_id"]
	if ok {
		igp := igprofile.GetIgProfile(igID.(string))
		status := igmedia.StatusShown
		if igp.Status != igprofile.StatusActive {
			status = igmedia.StatusHidden
		}
		return igmedia.UpdateStatusAll(igID.(string), status)
	}
	return "Param ig_id not found"
}
