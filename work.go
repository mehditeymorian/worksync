package worksync

import (
	"github.com/robfig/cron"
	"time"
)

type WorkStatus int8

const (
	// Queued Initial status of work before acquiring.
	Queued WorkStatus = iota + 1
	// Running indicates that the work is running by a worker.
	Running
	// Success status of work if the worker finish the work successfully.
	Success
	// Failed status of work if the worker finish the work with an error.
	Failed
)

type Work struct {
	// unique name for the work
	Name string
	// cron expression of when the worker is going to be executed
	Cron string
	// number of times the work must be executed.
	// Note that it is only guaranteed that
	// the number of workers acquiring this work not to surpass this value.
	MaxRun int64

	schedule     cron.Schedule
	previousTime time.Time
}

// parseSchedule parse cron expression as schedule.
func (w *Work) parseSchedule() error {
	schedule, err := cron.Parse(w.Cron)
	if err != nil {
		return err
	}

	w.schedule = schedule
	w.previousTime = time.Now()

	return nil
}
