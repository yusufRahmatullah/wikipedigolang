package igprofile

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"git.heroku.com/pg1-go-work/cmd/pg1-go/app/utils"

	"git.heroku.com/pg1-go-work/cmd/pg1-go/app/jobqueue"
)

var (
	commaVanisher = strings.NewReplacer(",", "")
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
	log.Println("==debug== run process")
	params := (*jq).Params
	igID, ok := params["ig_id"]
	if ok {
		page, err := utils.CreateWebPage()
		if err == nil {
			defer page.Close()
			log.Println("==debug== openning url")
			err := page.Open(fmt.Sprintf("https://www.instagram.com/%v", igID))
			if err == nil {
				log.Println("==debug== url openned")
				info, err := page.Evaluate(`function() {
					var folNode = document.querySelector("a[href$='followed_by_list']>span");
					var fol = folNode.getAttribute("title");
					return { followers: fol };
				}`)
				if err == nil {
					log.Printf("==debug== info: %v\n", info)
					data := info.(map[string]interface{})
					followers := commaVanisher.Replace(data["followers"].(string))
					log.Printf("==debug== followers text: %v\n", data["followers"])
					nbf, err := strconv.Atoi(followers)
					if err != nil {
						log.Printf("==debug== failed to convert str\n")
						return false
					}
					log.Printf("==debug== followers number: %v\n", nbf)
					return true
				}
				log.Println("==debug== failed to execute JS")
			} else {
				log.Println("==debug== failed to open url")
			}
		} else {
			log.Println("==debug== create page failed")
		}
	}
	return false
}
