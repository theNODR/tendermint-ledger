package helpers

import (
	"etl"

	"github.com/jasonlvhit/gocron"
)

type JobScheduler struct {
	isStarted		bool
	job				*etl.Job
	scheduler		*gocron.Scheduler
	schedulerStop	chan bool
}

func NewJobScheduler(job *etl.Job, typeName string, interval uint) *JobScheduler {
	scheduler := gocron.NewScheduler()
	scheduler.
		Every(uint64(interval)).
		Minutes().
		Do(etl.NewJobExecutor(job, typeName))

	return &JobScheduler{
		isStarted: false,
		job: job,
		scheduler: scheduler,
	}
}

func (j *JobScheduler) Start() {
	if j.isStarted {
		return
	}
	j.schedulerStop = j.scheduler.Start()
	j.job.Writer.Open()
	j.scheduler.RunAll()
}

func (j *JobScheduler) Stop() {
	if !j.isStarted {
		return
	}
	j.schedulerStop <- true
	j.job.Writer.Close()
	j.scheduler.RunAll()
}

func (j *JobScheduler) GetWriter() *etl.Writer{
	return j.job.Writer
}
