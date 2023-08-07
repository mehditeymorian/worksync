package main

import (
	"sync"
	"time"

	"github.com/mehditeymorian/worksync"
	"github.com/mehditeymorian/worksync/test/master/internal/config"
	"github.com/mehditeymorian/worksync/test/master/internal/db"
)

func main() {
	var wg sync.WaitGroup

	// load configs
	cfg := config.New("config.yaml")

	// open db connection
	database, err := db.NewConnection(cfg.Database)
	if err != nil {
		panic(err)
	}

	// opening a new sync db
	syncDB := worksync.NewDBConnection(database, "workers")

	// create list of works
	works := make([]*worksync.Work, 0)

	for _, job := range cfg.Jobs {
		works = append(works, &worksync.Work{
			Name:   job.Name,
			Cron:   job.Cron,
			MaxRun: job.MaxRun,
		})
	}

	// create a new scheduler
	scheduler, err := worksync.NewScheduler(syncDB, works, &worksync.SchedulerConfig{
		DurationBeforeSequence:    time.Duration(5) * time.Second,
		SchedulerCheckingInterval: "@every 2s",
		SequenceGenerator:         nil,
	})

	scheduler.StartSchedule()

	wg.Add(1)
	wg.Wait()
}
