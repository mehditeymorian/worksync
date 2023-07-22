package workersync

import (
	"github.com/robfig/cron"
	"time"
)

type WorkStatus int8

const (
	Queued WorkStatus = iota + 1
	Running
	Success
	Failed
)

type Work struct {
	Name   string
	Cron   string
	MaxRun uint

	schedule     cron.Schedule
	previousTime time.Time
}

func (w *Work) parseSchedule() {
	schedule, _ := cron.Parse(w.Cron)

	w.schedule = schedule
	w.previousTime = time.Now()
}
