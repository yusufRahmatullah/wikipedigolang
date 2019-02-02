package job

import (
	"fmt"
	"strings"

	"git.heroku.com/pg1-go-work/cmd/pg1-go/app/jobqueue"
	"git.heroku.com/pg1-go-work/cmd/pg1-go/app/logger"
	"git.heroku.com/pg1-go-work/cmd/pg1-go/app/utils"
)

const extJsFunc = `function() {
	var accNode = document.querySelectorAll("h2>a[title]");
	if (accNode.length >= 2) {
		var accCom = accNode[1];
		if (accCom) {
			var par = accCom.parentNode;
			if (par) {
				var gPar = par.parentNode;
				if (gPar.childElementCount >= 2) {
					var span = gPar.children[1];
					if (span) {
						return {text: span.innerText};
					}
				}
			}
		}
	}
	return {text: ""};
}`

var (
	extLogger = logger.NewLogger("PostExtractionJob", true, true)
)

// PostExtractionJob is the job for extracting all accounts
// that mentioned in a post
type PostExtractionJob struct{}

// NewPostExtractionJob instantiate PostExtractionJob instance
func NewPostExtractionJob() *PostExtractionJob {
	return &PostExtractionJob{}
}

// Name return PostExtractionJob name
func (job *PostExtractionJob) Name() string {
	return "PostExtractionJob"
}

func processExtractedAccounts(data map[string]interface{}) {
	extLogger.Debug("Processing post's text")
	text := data["text"].(string)
	words := strings.Split(text, " ")
	for _, word := range words {
		if strings.HasPrefix(word, "@") {
			acc := word[1:]
			cleanAccount := strings.Split(acc, "\n")[0]
			extLogger.Debug(fmt.Sprintf("Founc IG ID: %v", cleanAccount))
			jq := jobqueue.NewJobQueue("SingleAccountJob", map[string]interface{}{
				"ig_id": cleanAccount,
			})
			jobqueue.Save(jq)
		}
	}
}

func extractAccountFromPost(pID string) bool {
	wpw := utils.NewWebPageWrapper(extLogger)
	success := false
	if wpw != nil {
		defer wpw.Close()
		wpw.OnEvaluated(func(data map[string]interface{}) {
			processExtractedAccounts(data)
			success = true
		})
		wpw.OnError(func(err error) {
			success = false
		})
		wpw.OpenURL(fmt.Sprintf("https://www.instagram.com/p/%v", pID))
		wpw.Evaluate(extJsFunc)
	}
	return success
}

// Process extracts mentioned account in Ig Post
// returns true if process is success
func (job *PostExtractionJob) Process(jq *jobqueue.JobQueue) bool {
	extLogger.Debug("run post extraction process")
	params := jq.Params
	pID, ok := params["post_id"]
	if ok {
		return extractAccountFromPost(pID.(string))
	}
	extLogger.Info("Param post_id not found")
	return false
}
