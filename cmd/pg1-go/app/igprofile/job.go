package igprofile

import (
	"fmt"
	"strconv"
	"strings"

	"git.heroku.com/pg1-go-work/cmd/pg1-go/app/logger"

	"git.heroku.com/pg1-go-work/cmd/pg1-go/app/utils"

	"git.heroku.com/pg1-go-work/cmd/pg1-go/app/jobqueue"
)

var (
	commaVanisher = strings.NewReplacer(",", "")
	sajLogger     = logger.NewLogger("SingleAccountJob", true, true)
)

// SingleAccountJob is the job for crawling an account
// of IGProfile
type SingleAccountJob struct{}

// NewSingleAccountJob instantiate SingleAccountJob instance
// and assign its Name with "SingleAccountJob"
func NewSingleAccountJob() *SingleAccountJob {
	return &SingleAccountJob{}
}

// Name returns the SingleAccountJob
func (job *SingleAccountJob) Name() string {
	return "SingleAccountJob"
}

// Process executes job queue with the given params
// Returns true if process success
func (job *SingleAccountJob) Process(jq *jobqueue.JobQueue) bool {
	sajLogger.Debug("run process")
	params := (*jq).Params
	igID, ok := params["ig_id"]
	if ok {
		page, err := utils.CreateWebPage()
		if err == nil {
			defer page.Close()
			sajLogger.Debug("openning url")
			err := page.Open(fmt.Sprintf("https://www.instagram.com/%v", igID))
			if err == nil {
				sajLogger.Debug("url opened")
				info, err := page.Evaluate(`function() {
					var folNode = document.querySelector("a[href$='followed_by_list']>span");
					var folinNode = document.querySelector("a[href$='follows_list']>span");
					var postsNode = document.querySelector("a[href$='profile_posts']>span");
					var fol = folNode.getAttribute("title");
					var folin = folinNode.innerHTML;
					var posts = postsNode.innerHTML;
					var name = document.title.split("(")[0].trim();
					return { followers: fol, following: folin, posts: posts, name: name };
				}`)
				if err == nil {
					if info == nil {
						sajLogger.Fatal(fmt.Sprintf("info is nil, IG ID: %v", igID))
						return false
					}
					data := info.(map[string]interface{})
					followers := commaVanisher.Replace(data["followers"].(string))
					following := commaVanisher.Replace(data["following"].(string))
					posts := commaVanisher.Replace(data["posts"].(string))
					name := data["name"].(string)
					sajLogger.Debug(fmt.Sprintf("name: %v, followers: %v, following: %v, posts: %v", name, followers, following, posts))
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
					return Save(NewIgProfile(igID.(string), name, nbf, nbfin, np))
				}
				sajLogger.Fatal(fmt.Sprintf("Failed to execute JS on IG ID: %v", igID))
			} else {
				sajLogger.Fatal(fmt.Sprintf("Failed to open url on IG ID: %v", igID))
			}
		} else {
			sajLogger.Fatal("Failed to create page")
		}
	}
	return false
}
