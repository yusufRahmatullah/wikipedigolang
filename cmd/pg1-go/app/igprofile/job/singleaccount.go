package job

import (
	"fmt"
	"strconv"
	"strings"

	"git.heroku.com/pg1-go-work/cmd/pg1-go/app/igprofile"
	"git.heroku.com/pg1-go-work/cmd/pg1-go/app/jobqueue"
	"git.heroku.com/pg1-go-work/cmd/pg1-go/app/logger"
	"git.heroku.com/pg1-go-work/cmd/pg1-go/app/utils"
)

var (
	commaVanisher = strings.NewReplacer(",", "")
	sajLogger     = logger.NewLogger("SingleAccountJob", true, true)
	saJsFunction  = `function() {
		var folNode = document.querySelector("a[href$='followed_by_list']>span");
		var folinNode = document.querySelector("a[href$='follows_list']>span");
		var postsNode = document.querySelector("a[href$='profile_posts']>span");
		var ppNode = document.querySelector("span>img");
		var fol = folNode.getAttribute("title");
		var folin = folinNode.innerHTML;
		var posts = postsNode.innerHTML;
		if (ppNode) {
			var ppURL = ppNode.getAttribute("src");
		} else {
			var ppNode = document.querySelector("button>img");
			var ppURL = "";
			if (ppNode) {
				ppURL = ppNode.getAttribute("src");
			}
		}
		var name = document.title.split("(")[0].trim();
		return { followers: fol, following: folin, posts: posts, name: name, ppURL: ppURL };
	}`
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
	return igprofile.Save(igp)
}

func crawlIgID(igID string) bool {
	wpw := utils.NewWebPageWrapper(sajLogger)
	success := false
	if wpw != nil {
		defer wpw.Close()
		wpw.OnEvaluated(func(data map[string]interface{}) {
			success = processData(igID, data)
		})
		wpw.OnError(func(err error) {
			success = false
		})
		wpw.OpenURL(fmt.Sprintf("https://www.instagram.com/%v", igID))
		wpw.Evaluate(saJsFunction)
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
