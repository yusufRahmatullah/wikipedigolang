package job

import (
	"git.heroku.com/pg1-go-work/cmd/pg1-go/app/igprofile"
	"git.heroku.com/pg1-go-work/cmd/pg1-go/app/jobqueue"
	"git.heroku.com/pg1-go-work/cmd/pg1-go/app/logger"
)

var (
	snjLogger = logger.NewLogger("SearchNameJob", true, true)
)

// SearchNameJob is the job for get IG ID by search by name
type SearchNameJob struct{}

// NewSearchNameJob instantiate SearchNameJob instance
func NewSearchNameJob() *SearchNameJob {
	return &SearchNameJob{}
}

// Name returns SearchNameJob
func (job *SearchNameJob) Name() string {
	return "SearchNameJob"
}

// Process executes JobQueue with the given params
// Returns true if success
func (job *SearchNameJob) Process(jq *jobqueue.JobQueue) bool {
	params := jq.Params
	name, ok := params["ig_id"]
	if ok {
		igID := igprofile.FindByName(name.(string))
		if igID == "" {
			return false
		}
		jq := jobqueue.NewJobQueue("SingleAccountJob", map[string]interface{}{"ig_id": igID})
		return jobqueue.Save(jq)
	}
	snjLogger.Info("Param ig_id not found")
	return false
}
