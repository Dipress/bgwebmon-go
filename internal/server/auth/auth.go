package auth

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

const (
	invalidLoginOrPasswordMessage = "login or password is invalid"
	expiredTime                   = time.Hour * 24
)

// Authenticator interface for user authenticate
type Authenticator interface {
	Authenticate(*http.Request, *Response) error
}

type tokener interface {
	Token(userID int) (string, error)
}

type checker interface {
	Check(requestPassword, modelPassword string) error
}

type decoder interface {
	Decode(r *http.Request, model *Model) error
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

// ValidationErrorResponse struct
type ValidationErrorResponse struct {
	Status string `json:"status"`
	Fields []validationField
}

type validationField struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func (v ValidationErrorResponse) Error() string {
	return ""
}

// authenticator struct
type authenticator struct {
	db *sql.DB
	validator
	tokener
	checker
	decoder
}

// NewAuthenticator implements authenticator
func NewAuthenticator(db *sql.DB) Authenticator {
	return &authenticator{
		db:        db,
		validator: ozzoValidator{},
		tokener:   jwtToken{},
		checker:   passwordChecker{},
		decoder:   bodyDecoder{},
	}
}

// Authenticate implements interface Authenticator.Authenticate
func (u *authenticator) Authenticate(r *http.Request, resp *Response) error {
	request := Model{}
	if err := u.decoder.Decode(r, &request); err != nil {
		return err
	}

	if err := u.validator.Validate(request); err != nil {
		return err
	}

	user := Model{}
	if err := u.findByLogin(request.Login, &user); err != nil {
		return ErrorResponse{Status: "error", Message: invalidLoginOrPasswordMessage}
	}

	if err := u.checker.Check(request.Password, user.Password); err != nil {
		return err
	}

	token, err := u.tokener.Token(user.ID)
	if err != nil {
		return err
	}

	resp.Status = "ok"
	resp.Data = dataField{
		Token: token,
	}

	return nil
}

type passwordChecker struct{}

func (p passwordChecker) Check(requestPassword, modelPassword string) error {
	candidatePassword := passwordEncrypted(requestPassword)
	if candidatePassword != modelPassword {
		return ErrorResponse{Status: "error", Message: invalidLoginOrPasswordMessage}
	}
	return nil
}

// findByLogin find user by login in database
func (u *authenticator) findByLogin(login string, user *Model) error {
	row := u.db.QueryRow("SELECT id, login, pswd FROM user WHERE login = ?", login)
	switch err := row.Scan(&user.ID, &user.Login, &user.Password); err {
	case sql.ErrNoRows:
		return errors.Wrap(err, "login not found")
	default:
		errors.Wrap(err, "scan error")
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

type jwtToken struct{}

func (j jwtToken) Token(userID int) (string, error) {
	mySigningKey := []byte("secret")

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID": userID,
		"exp":    time.Now().Add(expiredTime).Unix(),
	})

	tokenString, err := token.SignedString(mySigningKey)
	if err != nil {
		return "", errors.Wrap(err, "unable to create token")
	}

	return tokenString, nil
}

type bodyDecoder struct{}

func (d bodyDecoder) Decode(r *http.Request, model *Model) error {
	err := json.NewDecoder(r.Body).Decode(&model)
	if err != nil {
		errors.Wrap(err, "unable to decode")
	}
	return nil
}
