package models

import (
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"log"
	"strings"
)

// User represents a user table in database.
type User struct {
	ID       int    `json:"id"`
	Login    string `json:"login"`
	Password string `json:"password"`
}

// FindByLogin find user by login in database
func FindByLogin(login string, db *sql.DB) (*User, error) {
	user := User{}
	rows, err := db.Query("select id, login, pswd from user where login = ?", login)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&user.ID, &user.Login, &user.Password)
		if err != nil {
			log.Fatal(err)
		}
	}

	return &user, nil
}

// PasswordEncrypted func take canditate password and return MD5 hex string
func PasswordEncrypted(password string) string {
	hash := md5.New()
	hash.Write([]byte(password))
	md := hash.Sum(nil)
	mdStr := hex.EncodeToString(md)
	return strings.ToUpper(mdStr)
}

// TODO: Stage 3 - create JWT token and give it a current user.

func tokenForUser(user *User) string {
	return ""
}
