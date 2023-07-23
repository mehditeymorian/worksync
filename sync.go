package workersync

type WorkResultFunc func()

type WorkSync struct {
	db *db
}

func New(db *db) *WorkSync {
	return &WorkSync{
		db: db,
	}
}

func (w *WorkSync) AcquireWork(workName, sequence string) (WorkResultFunc, WorkResultFunc, error) {
	id, err := w.db.acquireWork(workName, sequence)
	if err != nil {
		return nil, nil, err
	}

	success := func() {
		w.db.setWorkStatus(id, Success)
	}

	fail := func() {
		w.db.setWorkStatus(id, Failed)
	}

	return success, fail, nil
}
