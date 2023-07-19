package workersync

import (
	"context"
	"database/sql"
	"fmt"
)

type WorkResultFunc func()

type WorkerSyncer struct {
	db        *sql.DB
	tableName string
}

func New(db *sql.DB, tableName string) *WorkerSyncer {
	return &WorkerSyncer{
		db:        db,
		tableName: tableName,
	}
}

func (w *WorkerSyncer) AcquireWork(workName, sequence string) (WorkResultFunc, WorkResultFunc, error) {
	tx, err := w.db.BeginTx(context.Background(), &sql.TxOptions{
		Isolation: sql.LevelReadUncommitted,
	})
	if err != nil {
		return nil, nil, ErrFailToAcquireWork
	}

	_, err = tx.Exec("set @update_id := 0;")
	if err != nil {
		_ = tx.Rollback()

		return nil, nil, ErrFailToAcquireWork
	}

	_, err = tx.Query(`
	   update workers set status = ?, id = (select @update_id := id) where status = ? and name = ? and sequence_key = ? limit 1;
	`,
		Running,
		Queued,
		workName,
		sequence,
	)
	if err != nil {
		_ = tx.Rollback()

		return nil, nil, fmt.Errorf("%s: %w", ErrFailToAcquireWork.Error(), err)
	}

	result, err := tx.Query("select @update_id as id;")
	if err != nil {
		_ = tx.Rollback()

		return nil, nil, fmt.Errorf("%s: %w", ErrFailToAcquireWork.Error(), err)
	}

	columns, err := result.Columns()
	if err != nil {
		_ = tx.Rollback()

		return nil, nil, fmt.Errorf("%s: %w", ErrFailToAcquireWork.Error(), err)
	}

	if len(columns) == 0 {
		_ = tx.Rollback()

		return nil, nil, fmt.Errorf("%s: %w", ErrFailToAcquireWork.Error(), err)
	}

	if !result.Next() {
		_ = tx.Rollback()

		return nil, nil, fmt.Errorf("%s: %w", ErrFailToAcquireWork.Error(), err)
	}

	id := new(int64)

	err = result.Scan(id)
	if err != nil {
		_ = tx.Rollback()

		return nil, nil, fmt.Errorf("%s: %w", ErrFailToAcquireWork.Error(), err)
	}

	err = tx.Commit()
	if err != nil {
		_ = tx.Rollback()

		return nil, nil, fmt.Errorf("%s: %w", ErrFailToFinalizeAcquire.Error(), err)
	}

	success := func() {
		w.setWorkStatus(*id, Success)
	}

	fail := func() {
		w.setWorkStatus(*id, Failed)
	}

	return success, fail, nil
}

func (w *WorkerSyncer) setWorkStatus(id int64, status WorkStatus) error {
	result, err := w.db.Exec(`update workers set status = ?, updated_at = now() where id = ?;`,
		status,
		id)
	if err != nil {
		return fmt.Errorf("%s: %w", ErrFailToUpdateStatus, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("%s: %w", ErrFailToUpdateStatus, err)
	}

	if rowsAffected == 0 {
		return ErrFailToUpdateStatus
	}

	return nil
}
