package custom_errors

import "errors"

var (
	ErrUserNotFound    = errors.New("user with such email not found")
	ErrSessionNotFound = errors.New("session not found")
)

var (
	ErrEmptyFieldEmail    = errors.New("empty field email")
	ErrEmptyFieldUsername = errors.New("empty field username")
	ErrEmptyFieldPassword = errors.New("empty field password")
	ErrInvalidEmail       = errors.New("email sent is invalid")
)

var (
	ErrInvalidEmailOrPassword = errors.New("email or password is invalid")
	ErrJWTTimeExpired         = errors.New("JWT-token time expired")
)
