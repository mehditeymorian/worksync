package workersync

type WorkStatus int8

const (
	Queued WorkStatus = iota + 1
	Running
	Success
	Failed
)
