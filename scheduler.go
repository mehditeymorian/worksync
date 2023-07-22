package workersync

import (
	"fmt"
	"github.com/robfig/cron"
	"time"
)

type SequenceFunc func(schedule time.Time) string

type SchedulerConfig struct {
	DurationBeforeSequence    time.Duration
	SchedulerCheckingInterval string
	SequenceGenerator         SequenceFunc
}

var DefaultSchedulerConfig = &SchedulerConfig{
	DurationBeforeSequence:    time.Duration(10) * time.Minute,
	SchedulerCheckingInterval: "@every 1m",
	SequenceGenerator: func(schedule time.Time) string {
		return schedule.Format(time.DateTime)
	},
}

type WorkScheduler struct {
	cron      *cron.Cron
	works     []*Work
	entryChan chan *entry
}

func NewWorkScheduler(works []*Work, config *SchedulerConfig) (*WorkScheduler, error) {
	for _, work := range works {
		work.parseSchedule()
	}

	scheduler := WorkScheduler{
		works:     works,
		cron:      cron.New(),
		entryChan: make(chan *entry),
	}

	config = parseConfig(config)

	err := scheduler.cron.AddFunc(config.SchedulerCheckingInterval, func() {
		for _, work := range scheduler.works {
			nextSchedule := work.schedule.Next(work.previousTime)
			sub := nextSchedule.Sub(time.Now())
			if sub <= config.DurationBeforeSequence { // time to schedule
				scheduler.entryChan <- &entry{
					name:     work.Name,
					sequence: config.SequenceGenerator(nextSchedule),
					maxRun:   work.MaxRun,
				}
				work.previousTime = nextSchedule
			}
		}
	})
	if err != nil {
		return nil, fmt.Errorf("%s: %w", ErrFailToInitWorkScheduler.Error(), err)
	}

	go scheduler.workCreator()

	return &scheduler, nil
}

func (w *WorkScheduler) StartSchedule() {
	w.cron.Start()
}

func (w *WorkScheduler) StopSchedule() {
	w.cron.Stop()
}

func (w *WorkScheduler) workCreator() {
	for e := range w.entryChan {
		fmt.Printf("work: %s - sequence: %s - time: %s\n", e.name, e.sequence, time.Now().Format(time.TimeOnly))
	}
}

func parseConfig(config *SchedulerConfig) *SchedulerConfig {
	if config == nil {
		config = DefaultSchedulerConfig
		return config
	}

	if config.DurationBeforeSequence.Nanoseconds() == 0 {
		config.DurationBeforeSequence = DefaultSchedulerConfig.DurationBeforeSequence
	}

	if config.SchedulerCheckingInterval == "" {
		config.SchedulerCheckingInterval = DefaultSchedulerConfig.SchedulerCheckingInterval
	}

	if config.SequenceGenerator == nil {
		config.SequenceGenerator = DefaultSchedulerConfig.SequenceGenerator
	}

	return config
}
