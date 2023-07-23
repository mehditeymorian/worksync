package main

import (
	"database/sql"
	"github.com/mehditeymorian/worksync"
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	db, err := sql.Open("mysql", "root:@tcp(localhost:3306)/test?interpolateParams=false&parseTime=true&charset=utf8mb4")
	if err != nil {
		panic(err)
	}

	syncDB := worksync.NewDBConnection(db, "workers")

	name := "printer"
	cronExpr := "0,10,20,30,40,50 * * * * *"
	works := []*worksync.Work{
		{
			Name:   name,
			Cron:   cronExpr,
			MaxRun: 1,
		},
	}

	scheduler, err := worksync.NewScheduler(syncDB, works, &worksync.SchedulerConfig{
		DurationBeforeSequence:    time.Duration(5) * time.Second,
		SchedulerCheckingInterval: "@every 2s",
		SequenceGenerator:         nil,
	})

	scheduler.StartSchedule()

	var wg sync.WaitGroup
	wg.Add(1)
	wg.Wait()
}
