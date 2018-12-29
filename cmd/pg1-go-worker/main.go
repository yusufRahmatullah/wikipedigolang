package main

import (
	"log"
	"os"
	"strconv"
	"time"

	"git.heroku.com/pg1-go-work/cmd/pg1-go/app/utils"

	"git.heroku.com/pg1-go-work/cmd/pg1-go/app/jobqueue"
)

var (
	jobAssigner jobqueue.JobAssigner
)

func init() {
	// initialilze JobAssigner
	jobAssigner = jobqueue.NewJobAssigner()
	baseProcessor := jobqueue.NewBaseProcessor()
	jobAssigner.Register(baseProcessor)
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
	log.Println("==debug== [Consumer] running")
	for true {
		log.Println("==debug== [Consumer] retrieving queue")
		jobQueues = jobqueue.GetAll()
		if len(jobQueues) == 0 {
			log.Printf("==debug== [Consumer] Queue empty, waiting for %v seconds\n", waitTime)
			time.Sleep(time.Duration(waitTime) * time.Second)
		} else {
			log.Printf("==debug== [Consumer] Processing %v queues\n", len(jobQueues))
			for len(jobQueues) > 0 {
				jobQueue, jobQueues = jobQueues[0], jobQueues[1:]
				suc := jobAssigner.ProcessJobQueue(jobQueue)
				log.Printf("==debug== [Consumer] ProcessJob: %v\n", suc)
			}
		}
	}
}
