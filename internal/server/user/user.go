package user

import (
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	jwt "github.com/dgrijalva/jwt-go"

	"github.com/pkg/errors"
)

// Authenticator interface for user authenticate
type Authenticator interface {
	Authenticate(*http.Request, *Response) error
}

// Model represents a user table in database.
type Model struct {
	ID       int    `json:"id"`
	Login    string `json:"login"`
	Password string `json:"password"`
}

// Response represents response data
type Response struct {
	Status string    `json:"status"`
	Data   dataField `json:"data"`
}

type dataField struct {
	Token string `json:"token"`
}

// ErrorResponse represent response errors
type ErrorResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

// Implement the error interface
func (e ErrorResponse) Error() string {
	return e.Message
}

// authenticator struct
type authenticator struct {
	db *sql.DB
}

// NewAuthenticator implements authenticator
func NewAuthenticator(db *sql.DB) Authenticator {
	return &authenticator{
		db: db,
	}
}

// Authenticate implements interface Authenticator.Authenticate
func (u *authenticator) Authenticate(r *http.Request, resp *Response) error {
	request := Model{}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return errors.Wrap(err, "unable to decode")
	}

	user := Model{}
	if err := u.findByLogin(request.Login, &user); err != nil {
		return ErrorResponse{Status: "error", Message: "login or password a wrong"}

	}

	canditatePassword := passwordEncrypted(request.Password)

	if user.Password != canditatePassword {
		return ErrorResponse{Status: "error", Message: "login or password a wrong"}
	}

	token, err := createToken(user.ID)
	if err != nil {
		return errors.Wrap(err, "unable to create token")
	}

	resp.Status = "ok"
	resp.Data = dataField{
		Token: token,
	}

	return nil
}

// findByLogin find user by login in database
func (u *authenticator) findByLogin(login string, user *Model) error {
	row := u.db.QueryRow("SELECT id, login, pswd FROM user WHERE login = ?", login)
	switch err := row.Scan(&user.ID, &user.Login, &user.Password); err {
	case nil:
		return errors.Wrap(err, "scan error")
	case sql.ErrNoRows:
		return errors.Wrap(err, "login not found")
	default:
		errors.Wrap(err, "unprocessable entity")
	}
	return nil
}

// PasswordEncrypted func take canditate password and return MD5 hex string
func passwordEncrypted(password string) string {
	hash := md5.New()
	hash.Write([]byte(password))
	md := hash.Sum(nil)
	mdStr := hex.EncodeToString(md)
	return strings.ToUpper(mdStr)
}

func createToken(userID int) (string, error) {
	//The global secret key, don't show it, anybody.
	mySigningKey := []byte("secret")

	// Create new token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID": userID,
		"exp":    time.Now().Add(time.Hour * 24).Unix(),
	})

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString(mySigningKey)
	if err != nil {
		return "", errors.Wrap(err, "unable to create token")
	}

	return tokenString, nil
}
