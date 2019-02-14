package job

import (
	"git.heroku.com/pg1-go-work/cmd/pg1-go/app/igmedia"
	"git.heroku.com/pg1-go-work/cmd/pg1-go/app/igprofile"
	"git.heroku.com/pg1-go-work/cmd/pg1-go/app/jobqueue"
)

// MediaFromPostJob is the job for crawling all
// medias on a post
type MediaFromPostJob struct{}

// NewMediaFromPostJob instantiate MediaFromPostJob instance
func NewMediaFromPostJob() *MediaFromPostJob {
	return &MediaFromPostJob{}
}

// Name returns the MediaFromPostJob
func (job *MediaFromPostJob) Name() string {
	return "MediaFromPostJob"
}

// Process executes job queue with the given params
// Returns empty string if process success
func (job *MediaFromPostJob) Process(jq *jobqueue.JobQueue) string {
	params := jq.Params
	igID, ok := params["ig_id"]
	if !ok {
		return "Param ig_id not found"
	}
	postID, ok := params["post_id"]
	if !ok {
		return "Param post_id not found"
	}
	igms, err := igprofile.FetchMediaFromPost(igID.(string), postID.(string))
	if err != "" {
		return err
	}
	for _, igm := range igms {
		igmedia.Save(igm)
	}
	return ""
}
