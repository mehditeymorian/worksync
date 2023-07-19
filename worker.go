package workersync

type Worker struct {
	Name string
}

func NewWorker(name string) *Worker {
	return &Worker{
		Name: name,
	}
}
