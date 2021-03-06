package main

import (
	"os"
	"strconv"
	"time"

	igMediaJob "git.heroku.com/pg1-go-work/cmd/pg1-go/app/igmedia/job"

	"git.heroku.com/pg1-go-work/cmd/pg1-go/app/logger"

	igProfileJob "git.heroku.com/pg1-go-work/cmd/pg1-go/app/igprofile/job"
	"git.heroku.com/pg1-go-work/cmd/pg1-go/app/utils"

	"git.heroku.com/pg1-go-work/cmd/pg1-go/app/jobqueue"
)

var (
	jobAssigner *jobqueue.JobAssigner
	mainLogger  = logger.NewLogger("PG1-Go-Worker", false, true)
)

func init() {
	// initialize and migrate JobRank
	jobqueue.InitJobRank()
	jobqueue.MigrateJobRank()
	// initialilze JobAssigner
	jobAssigner = jobqueue.NewJobAssigner()
	singleAccountJob := igProfileJob.NewSingleAccountJob()
	updaterJob := igProfileJob.NewUpdaterJob()
	multiAccountJob := igProfileJob.NewMultiAccountJob()
	postExtractorJob := igProfileJob.NewPostExtractionJob()
	banAccountJob := igProfileJob.NewBanAccountJob()
	searchNameJob := igProfileJob.NewSearchNameJob()
	updateStatusJob := igMediaJob.NewUpdateIgMediaStatusJob()
	accFromPostJob := igProfileJob.NewAccountFromPostJob()
	mediaFromPostJob := igMediaJob.NewMediaFromPostJob()
	singleUpdaterJob := igProfileJob.NewSingleUpdaterJob()
	jobAssigner.Register(singleAccountJob)
	jobAssigner.Register(updaterJob)
	jobAssigner.Register(multiAccountJob)
	jobAssigner.Register(postExtractorJob)
	jobAssigner.Register(banAccountJob)
	jobAssigner.Register(searchNameJob)
	jobAssigner.Register(updateStatusJob)
	jobAssigner.Register(accFromPostJob)
	jobAssigner.Register(mediaFromPostJob)
	jobAssigner.Register(singleUpdaterJob)
}

func getWaitTime() int {
	waitTimeStr := os.Getenv("WAIT_TIME")
	waitTime, err := strconv.Atoi(waitTimeStr)
	if err != nil {
		mainLogger.Warning("$WAIT_TIME not found, use default")
		waitTime = 5
	}
	return waitTime
}

func consumeJobs(waitTime int) {
	var jobQueues []jobqueue.JobQueue
	var jobQueue jobqueue.JobQueue

	for true {
		jobQueues = jobqueue.GetAll()
		if len(jobQueues) == 0 {
			time.Sleep(time.Duration(waitTime) * time.Second)
		} else {
			for len(jobQueues) > 0 {
				jobQueue, jobQueues = jobQueues[0], jobQueues[1:]
				jobAssigner.ProcessJobQueue(&jobQueue)
			}
		}
	}
}

func main() {
	// open phantom JS
	err := utils.DefaultProcess.Open()
	defer utils.DefaultProcess.Close()
	if err != nil {
		mainLogger.Error("Failed to open phantomjs process", err)
	}
	consumeJobs(getWaitTime())
}
