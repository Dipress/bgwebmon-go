package server

import (
	"database/sql"
	"net/http"

	"github.com/dipress/bgwebmon-go/internal/server/user"
	"github.com/julienschmidt/httprouter"
)

//New func create server
func New(db *sql.DB) *http.Server {
	mux := httprouter.New()
	mux.POST("/auth", postLogin(user.NewAuthenticator(db)))

	return &http.Server{
		Addr:    "127.0.0.1:5000",
		Handler: mux,
	}
}
