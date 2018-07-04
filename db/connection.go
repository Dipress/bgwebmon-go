package db

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
	"github.com/pkg/errors"
)

var db *sql.DB

// OpenDB return a sql connection instance or error
func OpenDB(driver, dsn string) (*sql.DB, error) {
	db, err := sql.Open(driver, dsn)
	if err != nil {
		return nil, errors.Wrap(err, "open sql connection failed: %v\n")
	}
	if err := db.Ping(); err != nil {
		return nil, errors.Wrap(err, "mysql ping failure: %v\n")
	}
	defer db.Close()

	return db, nil

}
