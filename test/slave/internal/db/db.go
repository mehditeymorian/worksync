package db

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

func NewConnection(cfg Config) (*sql.DB, error) {
	conn, err := sql.Open("mysql", cfg.DNS())
	if err != nil {
		return nil, err
	}

	return conn, nil
}
