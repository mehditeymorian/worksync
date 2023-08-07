package db

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

func NewConnection(cfg Config) (*sql.DB, error) {
	conn, err := sql.Open("mysql", "root:@tcp(localhost:3306)/test?interpolateParams=false&parseTime=true&charset=utf8mb4")
	if err != nil {
		return nil, err
	}

	return conn, nil
}
