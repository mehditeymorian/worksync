package workersync

type JobStatus int8

const (
	Running JobStatus = iota + 1
	Success
	Failed
)

type Job struct {
	Name      string
	Type      uint64
	Timestamp string
	Status    JobStatus
}
