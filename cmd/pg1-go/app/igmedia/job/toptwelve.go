package job

import (
	"strings"

	"git.heroku.com/pg1-go-work/cmd/pg1-go/app/igmedia"
	"git.heroku.com/pg1-go-work/cmd/pg1-go/app/igprofile"
	"git.heroku.com/pg1-go-work/cmd/pg1-go/app/jobqueue"
	"git.heroku.com/pg1-go-work/cmd/pg1-go/app/logger"
)

var (
	ttLogger = logger.NewLogger("TopTwelveJob", true, true)
)

// TopTwelveJob is the job for get top twelve post of IgProfile
type TopTwelveJob struct{}

// NewTopTwelveJob instantiate TopTwelveJob instance
func NewTopTwelveJob() *TopTwelveJob {
	return &TopTwelveJob{}
}

// Name returns TopTwelveJob
func (job *TopTwelveJob) Name() string {
	return "TopTwelveJob"
}

// Process executes JobQueue with the given params
// Returns empty string if success
func (job *TopTwelveJob) Process(jq *jobqueue.JobQueue) string {
	params := jq.Params
	igID, ok := params["ig_id"]
	if ok {
		topTwelve, errMsg := igprofile.TopTwelveMedia(igID.(string))
		if errMsg != "" {
			return errMsg
		}
		for _, node := range topTwelve {
			if strings.Contains(node.AccessibilityCaption, "people") || strings.Contains(node.AccessibilityCaption, "person") {
				igmedia.Save(igmedia.NewIgMedia(node.ID, igID.(string), node.DisplayURL))
			}
		}
		return ""
	}
	ttLogger.Info("Param ig_id not found")
	return "Param ig_id not found"
}
