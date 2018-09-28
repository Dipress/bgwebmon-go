package server

import (
	"database/sql"
	"net/http"

	"github.com/dipress/bgwebmon-go/internal/server/auth"
	"github.com/julienschmidt/httprouter"
	"github.com/rileyr/middleware"
)

func setHeader(fn httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		w.Header().Set("Content-type", "application/json")
		fn(w, r, p)
	}
}

// New func create server
func New(db *sql.DB, secretKey string) *http.Server {
	mux := httprouter.New()

	s := middleware.NewStack()
	s.Use(setHeader)

	mux.POST("/auth", s.Wrap(postLogin(auth.NewAuthenticator(db, secretKey))))

	return &http.Server{
		Addr:    "127.0.0.1:5000",
		Handler: mux,
	}
}
