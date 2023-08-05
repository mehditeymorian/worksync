<h1 align="center">
<img alt="Koi logo" src="assets/image.png" width="500px"/><br/>
Work Synchronizer
</h1>
<p align="center">Schedule works for different types of workers on different machines backed by only a single `MySQL` table.
</p>

<p align="center">
<a href="https://pkg.go.dev/github.com/mehditeymorian/worksync?tab=doc"target="_blank">
    <img src="https://img.shields.io/badge/Go-1.20+-00ADD8?style=for-the-badge&logo=go" alt="go version" />
</a>&nbsp;
<img src="https://img.shields.io/badge/license-MIT-blue?style=for-the-badge&logo=none" alt="license" />

<img src="https://img.shields.io/badge/Version-0.0.1-informational?style=for-the-badge&logo=none" alt="version" />
</p>

## Why this library
- Small Codebase
- Minimum dependency on other libraries
- Flexibility to integrate into workers

## Download
```
go get github.com/mehditeymorian/worksync@latest
```

## How the Sync Happens?
All the works are managed by a `scheduler` backed by `MYSQL`. The scheduler takes a list of work blueprints with the following attributes:
```go
type Work struct {
	// unique name for the work
	Name   string
	// cron expression of when the worker is going to be executed
	Cron   string
	// number of times the work must be executed.
	// Note that it is only guaranteed that 
	// the number of workers acquiring this work not to surpass this value.
	MaxRun int64
}
```
The Scheduler then creates works before the schedule of each worker based on the `MaxRun` value. Then at the scheduled time, each worker tries to acquire work by updating the work status. If there is work with a `Queued` status, the worker will acquire the work. otherwise, it fails.
> Note: each sequence of works is known by the worker `name` and a `sequence` value. These two values must be the same on the scheduler and on the worker side.
> 
> By default, the scheduler set `sequence` as the execution time formatted by `time.DateTime`. The sequence generator can be customized in the schedule's config.

## How to Use
Create the following table in MYSQL
```sql
create table if not exists [DATABASE].[TABLE]
(
    id          int auto_increment
        primary key,
    created_at  datetime     null,
    name        varchar(255) null,
    status      int8         null,
    sequence    varchar(255) null,
    started_at  datetime     null,
    finished_at datetime     null
);
```
Then connect as follows:
```go
package main

import (
	...
        // import the proper driver for connection
        _ "github.com/go-sql-driver/mysql"
)

func main() {
    db, err := sql.Open("mysql", dsn)
    if err != nil {
        panic(err)
    }
    
    syncDB := worksync.NewDBConnection(db, TableName)

}
```


### Worker Side
On the client-side, it only requires calling the `acquireWork` before doing the work.
```go
package main

func main() {

	workSync := worksync.NewSyncer(syncDB)
	
	fancyWorker.do(func() {
		success, fail, err := workSync.AcquireWork(name, now.Format(time.DateTime))
		if err != nil {
                    return
		}
		
		// actual work
	})

}
```

### Manager Side
```go
package main

import (
	"worksync"
)

func main() {

	works := []*worksync.Work{
		{
			Name:   name,
			Cron:   cronExpr,
			MaxRun: maxRun,
		},
		// list of works ...
	}

	scheduler, err := worksync.NewScheduler(syncDB, works, &worksync.SchedulerConfig{
		// when to create the works. in this case it will create the work if duration before execution time is less than 5sec.
		DurationBeforeSequence:    time.Duration(5) * time.Second,
		// interval to check for creating works
		SchedulerCheckingInterval: "@every 2s",
		SequenceGenerator:         nil,
	})

	// non-blocking
	scheduler.StartSchedule()


}
```


## Cron Expressions
We use [robfig/cron](https://github.com/robfig/cron) for handling Crons. Please refer to their documentation on how to write cron expressions. [Link](https://pkg.go.dev/github.com/robfig/cron)
