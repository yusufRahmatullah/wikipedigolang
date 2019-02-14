package job

import (
	"git.heroku.com/pg1-go-work/cmd/pg1-go/app/igprofile"
	"git.heroku.com/pg1-go-work/cmd/pg1-go/app/jobqueue"
)

// AccountFromPostJob is the job for crawling an account
// on a post
type AccountFromPostJob struct{}

// NewAccountFromPostJob instantiate AccountFromPostJob instance
func NewAccountFromPostJob() *AccountFromPostJob {
	return &AccountFromPostJob{}
}

// Name returns the AccountFromPostJob
func (job *AccountFromPostJob) Name() string {
	return "AccountFromPostJob"
}

// Process executes job queue with the given params
// Returns empty string if process success
func (job *AccountFromPostJob) Process(jq *jobqueue.JobQueue) string {
	params := jq.Params
	postID, ok := params["post_id"]
	if ok {
		accs, err := igprofile.FetchAccountFromPost(postID.(string))
		if err != "" {
			return err
		}
		for _, acc := range accs {
			jq := jobqueue.NewJobQueue(
				"SingleAccountJob",
				map[string]interface{}{"ig_id": acc},
			)
			jobqueue.Save(jq)
		}
		return ""
	}
	return "Param post_id not found"
}
