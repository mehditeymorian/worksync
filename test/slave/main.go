package main

import (
	"fmt"
	"log"
	"time"

	"github.com/mehditeymorian/worksync"
	"github.com/mehditeymorian/worksync/test/slave/internal/config"
	"github.com/mehditeymorian/worksync/test/slave/internal/db"

	"github.com/robfig/cron"
)

func main() {
	channel := make(chan struct{})

	// load configs
	cfg := config.New("config.yaml")

	// open db connection
	database, err := db.NewConnection(cfg.Database)
	if err != nil {
		panic(err)
	}

	// creating a sync db and a syncer
	syncDB := worksync.NewDBConnection(database, "workers")
	workSync := worksync.NewSyncer(syncDB)

	// creating cron workers
	workers := cron.New()

	// adding a new worker
	_ = workers.AddFunc(cfg.Job.Cron, func() {
		success, _, er := workSync.AcquireWork(cfg.Job.Name, time.Now().Format(time.DateTime))
		if er != nil {
			log.Printf("[cron: %s][job: %s] didn't acquire the work error=%v\n", cfg.Job.Cron, cfg.Job.Name, err)

			return
		}

		log.Println(fmt.Sprintf("[cron: %s][job: %s] done!", cfg.Job.Cron, cfg.Job.Name))

		success()
	})

	workers.Start()

	log.Println(fmt.Sprintf("slave [%s] is set.", cfg.Name))

	<-channel
}
