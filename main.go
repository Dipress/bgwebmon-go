package main

import (
	"database/sql"
	"flag"
	"log"

	"github.com/dipress/bgwebmon-go/internal/server"
	_ "github.com/go-sql-driver/mysql"
	"github.com/pkg/errors"
)

var (
	dsn = flag.String("dsn", "username:password@tcp(127.0.0.1:3306)/bgbilling_development",
		"mysql connection string")
	secretKey = flag.String("secret-key", "", "secret key for jwt token")
)

func main() {
	flag.Parse()

	db, err := openDB(*dsn)
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	server := server.New(db, *secretKey)
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}

// OpenDB return sql instance
func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, errors.Wrap(err, "open sql connection failed: %v\n")
	}
	if err := db.Ping(); err != nil {
		return nil, errors.Wrap(err, "mysql ping failure: %v\n")
	}

	return db, nil

}
