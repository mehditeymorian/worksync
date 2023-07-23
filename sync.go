package worksync

type WorkResultFunc func()

type WorkSync struct {
	db *db
}

func NewSyncer(db *db) *WorkSync {
	return &WorkSync{
		db: db,
	}
}

// AcquireWork check database and try to acquire a work.
// It returns an error if there is no work to acquire or fail in the process of acquiring a work.
// If the acquiring process is successful, it returns success and fail functions to determine the result of the work.
// Calling either of these functions will update the status of the work in the database.
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
