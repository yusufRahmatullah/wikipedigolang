package job

import (
	"fmt"
	"strconv"
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

func processData(igID string, data map[string]interface{}) bool {
	followers := commaVanisher.Replace(data["followers"].(string))
	following := commaVanisher.Replace(data["following"].(string))
	posts := commaVanisher.Replace(data["posts"].(string))
	name := data["name"].(string)
	ppURL := data["ppURL"].(string)
	nbf, err := strconv.Atoi(followers)
	if err != nil {
		sajLogger.Fatal(fmt.Sprintf("Failed to convert follower text to int: %v", followers))
		return false
	}
	sajLogger.Debug(fmt.Sprintf("Followers: %v", nbf))

	nbfin, err := strconv.Atoi(following)
	if err != nil {
		sajLogger.Fatal(fmt.Sprintf("Failed to convert following text: %v", following))
		return false
	}
	sajLogger.Debug(fmt.Sprintf("Following: %v", nbfin))

	np, err := strconv.Atoi(posts)
	if err != nil {
		sajLogger.Fatal(fmt.Sprintf("Failed to convert posts text: %v", posts))
		return false
	}
	sajLogger.Debug(fmt.Sprintf("Posts: %v", np))
	igpBuilder := igprofile.NewBuilder()
	igpBuilder = igpBuilder.SetIGID(igID).SetName(name)
	igpBuilder = igpBuilder.SetFollowers(nbf).SetFollowing(nbfin)
	igp := igpBuilder.SetPosts(np).SetPpURL(ppURL).Build()
	return igprofile.SaveOrUpdate(igp)
}

func crawlIgID(igID string) bool {
	igp := igprofile.FetchIgProfile(igID)
	success := false
	if igp != nil {
		success = igprofile.SaveOrUpdate(igp)
	}
	return success
}

// Process executes job queue with the given params
// Returns true if process success
func (job *SingleAccountJob) Process(jq *jobqueue.JobQueue) bool {
	sajLogger.Debug("run process")
	params := (*jq).Params
	igID, ok := params["ig_id"]
	if ok {
		return crawlIgID(igID.(string))
	}
	return false
}
