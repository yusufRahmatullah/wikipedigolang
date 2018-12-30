package main

import (
	"log"
	"os"
	"strconv"
	"time"

	"git.heroku.com/pg1-go-work/cmd/pg1-go/app/igprofile"
	"git.heroku.com/pg1-go-work/cmd/pg1-go/app/utils"

	"git.heroku.com/pg1-go-work/cmd/pg1-go/app/jobqueue"
)

var (
	jobAssigner *jobqueue.JobAssigner
)

func init() {
	// initialilze JobAssigner
	jobAssigner = jobqueue.NewJobAssigner()
	singleAccountJob := igprofile.NewSingleAccountJob()
	jobAssigner.Register(singleAccountJob)
}

func main() {
	var jobQueues []jobqueue.JobQueue
	var jobQueue jobqueue.JobQueue
	waitTimeStr := os.Getenv("WAIT_TIME")
	waitTime, err := strconv.Atoi(waitTimeStr)

	utils.HandleError(
		err,
		"$WAIT_TIME not found. Use default instead",
		func() {
			waitTime = 5
		},
	)

	err = utils.DefaultProcess.Open()
	defer utils.DefaultProcess.Close()
	if err != nil {
		log.Println("==debug== failed open phantomjs process")
		panic(err)
	}

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
