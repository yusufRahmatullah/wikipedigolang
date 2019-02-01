package job

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"git.heroku.com/pg1-go-work/cmd/pg1-go/app/igprofile"

	"git.heroku.com/pg1-go-work/cmd/pg1-go/app/jobqueue"
	"git.heroku.com/pg1-go-work/cmd/pg1-go/app/logger"
	"git.heroku.com/pg1-go-work/cmd/pg1-go/app/utils"
)

const maJsAsyncFuncTemp = `function() {
    window.retData = new Set();
    function run(times) {
        if (times > 0) {
			window.scroll(0, document.body.scrollHeight);
			var pNodes = document.querySelectorAll("a[href^='/p/']");
			for (var i=0, n=pNodes.length; i<n; i++) {
				var href = pNodes[i].getAttribute("href");
				window.retData.add(href);
			}
			setTimeout(run.bind(null, times-1), 2000);
		}
    }
    setTimeout(run.bind(null, %d), 1);
}`

const maJsSyncFunc = `function() {
	return {posts: Array.from(window.retData)};
}`

const defaultAsyncLoop = 10

var (
	asyncLoop     int
	asyncWaitTime int
	maJsAsyncFunc string
	majLogger     = logger.NewLogger("MultiAccountJob", true, true)
)

func init() {
	var err error
	alStr := os.Getenv("ASYNC_LOOP")
	if alStr == "" {
		majLogger.Info("$ASYNC_LOOP not found, use default")
		asyncLoop = defaultAsyncLoop
	} else {
		asyncLoop, err = strconv.Atoi(alStr)
		if err != nil {
			majLogger.Info(fmt.Sprintf("Failed to convert %v to int", alStr))
			asyncLoop = defaultAsyncLoop
		}
	}
	asyncWaitTime = asyncLoop*2 + 5
	maJsAsyncFunc = fmt.Sprintf(maJsAsyncFuncTemp, asyncWaitTime)
}

// MultiAccountJob is the job for crawling an account
// that contains multiple accounts in its posts
type MultiAccountJob struct{}

// NewMultiAccountJob instantiate MultiAccountJob instance
func NewMultiAccountJob() *MultiAccountJob {
	return &MultiAccountJob{}
}

// Name returns the MultiAccountJob
func (job *MultiAccountJob) Name() string {
	return "MultiAccountJob"
}

func processPostsData(data map[string]interface{}) {
	majLogger.Debug("Processing posts data")
	posts := data["posts"].([]interface{})
	majLogger.Debug(fmt.Sprintf("There are %v posts", len(posts)))
	for _, post := range posts {
		ss := strings.Split(post.(string)[3:], "/")
		pID := ss[0]
		if pID != "" {
			jq := jobqueue.NewJobQueue("PostExtractionJob", map[string]interface{}{
				"post_id": pID,
			})
			jobqueue.Save(jq)
		}
	}
}

func crawlMultiIgID(igID string) {
	wpw := utils.NewWebPageWrapper(majLogger)
	if wpw != nil {
		defer wpw.Close()
		wpw.OnAsyncEvaluated(func() {
			majLogger.Debug(fmt.Sprintf("Async Func success. Waiting for %v seconds", asyncWaitTime))
			time.Sleep(time.Duration(asyncWaitTime) * time.Second)
			wpw.Evaluate(maJsSyncFunc)
		})
		wpw.OnEvaluated(func(data map[string]interface{}) {
			processPostsData(data)
		})
		wpw.OpenURL(fmt.Sprintf("https://www.instagram.com/%v", igID))
		wpw.EvaluateAsync(maJsAsyncFunc)
	}
}

// Process executes JobQueue with the given params
// Always returns true due the process is async
func (job *MultiAccountJob) Process(jq *jobqueue.JobQueue) bool {
	majLogger.Debug("run process")
	params := (*jq).Params
	igID, ok := params["ig_id"]
	if ok {
		cleanID := igprofile.CleanIgIDParams(igID.(string))
		suc := igprofile.SaveMultiAcc(cleanID)
		majLogger.Info(fmt.Sprintf("Save multi account: %v", suc))
		crawlMultiIgID(cleanID)
	} else {
		majLogger.Info("Param ig_id not found")
	}
	return true
}
