package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/dipress/crmifc_manager/models"
	"github.com/julienschmidt/httprouter"
)

// PostLogin func return login page
func PostLogin(db *sql.DB) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		userRequest := models.User{}

		decoder := json.NewDecoder(r.Body)

		err := decoder.Decode(&userRequest)
		if err != nil {
			log.Fatal(err)
			w.Write([]byte("Unprocesable Entity"))
			return
		}

		user, err := models.FindByLogin(userRequest.Login, db)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Unauthorized, wrong login"))
			return
		}

		canditatePassword := models.PasswordEncrypted(userRequest.Password)

		if user.Password != canditatePassword {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Unauthorized, wrong password"))
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(strconv.Itoa(user.ID)))
		return
	}
}
