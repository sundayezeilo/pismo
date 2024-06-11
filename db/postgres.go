package db

import (
	"database/sql"

	_ "github.com/lib/pq"
)

func NewPostgresDB(driver string, connStr string) (*sql.DB, error) {
	db, err := sql.Open(driver, connStr)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
