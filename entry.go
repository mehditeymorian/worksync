package workersync

type entry struct {
	name     string
	sequence string
	maxRun   int64
}
