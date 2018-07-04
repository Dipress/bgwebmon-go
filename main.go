package main

import (
	"database/sql"
	"log"

	"net/http"

	"github.com/dipress/crmifc_manager/handlers"
	_ "github.com/go-sql-driver/mysql"
	"github.com/julienschmidt/httprouter"
)

// Database struct
type Database struct {
	db *sql.DB
}

// NewDatabase initialize Database
func NewDatabase(db *sql.DB) *Database {
	return &Database{db: db}
}

func main() {

	db, err := sql.Open("mysql", "username:password@tcp(127.0.0.1:3306)/bgbilling_development")
	if err != nil {
		log.Fatal(err)
	}

	router := httprouter.New()
	router.POST("/login", handlers.PostLogin(db))
	router.GET("/appeals", handlers.GetAppeal)
	router.GET("/claims", handlers.GetClaim)

	log.Fatal(http.ListenAndServe(":5000", router))

}
