package main

import (
	"database/sql"
	"fmt"
	"github.com/mehditeymorian/worksync"
	"github.com/robfig/cron"
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

	workSync := worksync.NewSyncer(syncDB)

	workers := cron.New()

	n := 5

	for i := 0; i < n; i++ {
		workers.AddFunc(cronExpr, func() {
			now := time.Now()

			success, _, err := workSync.AcquireWork(name, now.Format(time.DateTime))
			if err != nil {
				fmt.Println("didn't acquire the work", err)
				return
			}

			<-time.After(time.Second)
			fmt.Println("hello world", now.Format(time.DateTime))
			success()
		})
	}

	workers.Start()

	var wg sync.WaitGroup
	wg.Add(1)
	wg.Wait()
}
