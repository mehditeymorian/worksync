package worksync

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"math"
	"strings"
	"text/template"
)

const DefaultTableName = "works"

type db struct {
	db *sql.DB

	queryParserData map[string]any
}

func NewDBConnection(conn *sql.DB, tableName string) *db {
	if tableName == "" {
		tableName = DefaultTableName
	}

	return &db{
		db: conn,

		queryParserData: map[string]any{
			"Table": tableName,
		},
	}
}

func (d *db) createWork(e *entry) error {
	tx, err := d.db.BeginTx(context.Background(), &sql.TxOptions{
		Isolation: sql.LevelSerializable,
	})
	if err != nil {
		return err
	}

	rows, err := tx.Query(d.prepareQuery("select count(*) from {{.Table}} where name = ? and sequence = ?"),
		e.name,
		e.sequence)
	if err != nil {
		_ = tx.Rollback()

		return err
	}

	count, err := d.count(rows)
	if err != nil {
		_ = tx.Rollback()

		return err
	}

	rows.Close()

	createCount := int(math.Max(0, float64(e.maxRun-count)))

	if createCount > 0 {
		builder := new(strings.Builder)
		builder.WriteString(d.prepareQuery("insert into {{.Table}} (created_at, name, status, sequence) values "))
		for i := 0; i < createCount; i++ {
			builder.WriteString(fmt.Sprintf("(now(), '%s', %d, '%s')", e.name, Queued, e.sequence))
			if i+1 == createCount {
				builder.WriteByte(';')
			} else {
				builder.WriteByte(',')
			}
		}

		_, err := tx.Exec(builder.String())
		if err != nil {
			_ = tx.Rollback()

			return err
		}

		err = tx.Commit()
		if err != nil {
			return err
		}
	}

	return nil
}

func (d *db) acquireWork(workName, sequence string) (int64, error) {
	tx, err := d.db.BeginTx(context.Background(), &sql.TxOptions{
		Isolation: sql.LevelReadUncommitted,
	})
	if err != nil {
		return 0, ErrFailToAcquireWork
	}

	_, err = tx.Exec("set @update_id := 0;")
	if err != nil {
		_ = tx.Rollback()

		return 0, ErrFailToAcquireWork
	}

	_, err = tx.Query(d.prepareQuery(`
	   update {{.Table}} set status = ?, started_at = now(), id = (select @update_id := id) where status = ? and name = ? and sequence = ? limit 1;
	`),
		Running,
		Queued,
		workName,
		sequence,
	)
	if err != nil {
		_ = tx.Rollback()

		return 0, fmt.Errorf("%s: %w", ErrFailToAcquireWork.Error(), err)
	}

	result, err := tx.Query("select @update_id as id;")
	if err != nil {
		_ = tx.Rollback()

		return 0, fmt.Errorf("%s: %w", ErrFailToAcquireWork.Error(), err)
	}

	columns, err := result.Columns()
	if err != nil {
		_ = tx.Rollback()

		return 0, fmt.Errorf("%s: %w", ErrFailToAcquireWork.Error(), err)
	}

	if len(columns) == 0 {
		_ = tx.Rollback()

		return 0, fmt.Errorf("%s: %w", ErrFailToAcquireWork.Error(), err)
	}

	if !result.Next() {
		_ = tx.Rollback()

		return 0, fmt.Errorf("%s: %w", ErrFailToAcquireWork.Error(), err)
	}

	id := new(int64)

	err = result.Scan(id)
	if err != nil {
		_ = tx.Rollback()

		return 0, fmt.Errorf("%s: %w", ErrFailToAcquireWork.Error(), err)
	}

	err = tx.Commit()
	if err != nil {
		_ = tx.Rollback()

		return 0, fmt.Errorf("%s: %w", ErrFailToFinalizeAcquire.Error(), err)
	}

	return *id, nil
}

func (d *db) setWorkStatus(id int64, status WorkStatus) error {
	result, err := d.db.Exec(`update workers set status = ?, finished_at = now() where id = ?;`,
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

func (d *db) count(result *sql.Rows) (int64, error) {
	columns, err := result.Columns()
	if err != nil {
		return 0, err
	}

	if len(columns) == 0 {
		return 0, errors.New("no result is returned from the query")
	}

	if !result.Next() {
		return 0, errors.New("no data is available from the query result")

	}

	count := new(int64)

	err = result.Scan(count)
	if err != nil {
		return 0, err
	}

	return *count, nil
}

func (d *db) prepareQuery(query string) string {
	temp := template.Must(template.New("query_parser").Parse(query))

	result := new(strings.Builder)

	err := temp.Execute(result, d.queryParserData)
	if err != nil {
		return ""
	}

	return result.String()
}
