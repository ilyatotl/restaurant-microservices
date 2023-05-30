package core

import (
	"authorization/internal/app/custom_errors"
	"authorization/internal/app/sessions"
	"authorization/internal/app/user"
	"context"
	"github.com/golang-jwt/jwt"
	password "github.com/vzglad-smerti/password_hash"
	"net/mail"
	"time"
)

var secretKey = []byte("AzAwCIxizF4cIgMLk4VrWhLYFZJlnOS52EF63K2I4OM=")

type Core struct {
	UsersRepo    *user.UsersRepository
	SessionsRepo *sessions.SessionRepository
}

func NewCore(userRepo *user.UsersRepository, sessionRepo *sessions.SessionRepository) *Core {
	return &Core{
		UsersRepo:    userRepo,
		SessionsRepo: sessionRepo,
	}
}

func (c *Core) validMailAddress(email string) bool {
	if _, err := mail.ParseAddress(email); err != nil {
		return false
	}
	return true
}

func (c *Core) RegisterUser(ctx context.Context, u *user.UserModel) (int64, error) {
	if u.Email == "" {
		return 0, custom_errors.ErrEmptyFieldEmail
	}
	if u.Username == "" {
		return 0, custom_errors.ErrEmptyFieldUsername
	}
	if u.Password == "" {
		return 0, custom_errors.ErrEmptyFieldPassword
	}
	if !c.validMailAddress(u.Email) {
		return 0, custom_errors.ErrInvalidEmail
	}

	encrypted_pas, err := password.Hash(u.Password)
	if err != nil {
		return 0, err
	}
	u.PasswordHash = encrypted_pas

	return c.UsersRepo.Add(ctx, u)
}

func (c *Core) AuthorizeUser(ctx context.Context, email string, p string) (string, error) {
	if email == "" {
		return "", custom_errors.ErrEmptyFieldEmail
	}
	if p == "" {
		return "", custom_errors.ErrEmptyFieldPassword
	}

	u, err := c.UsersRepo.Get(ctx, email)
	if err != nil {
		return "", err
	}

	if ok, err := password.Verify(u.PasswordHash, p); err != nil {
		return "", err
	} else if !ok {
		return "", custom_errors.ErrInvalidEmailOrPassword
	}

	exp := time.Now().Add(60 * time.Minute)
	id, err := c.SessionsRepo.Add(ctx, &sessions.SessionModel{
		UserID:    u.ID,
		ExpiresAt: exp,
	})
	if err != nil {
		return "", err
	}

	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["exp"] = exp
	claims["session_id"] = id

	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}
	c.SessionsRepo.Update(ctx, &sessions.SessionModel{
		ID:           id,
		SessionToken: tokenString,
	})

	return tokenString, nil
}

func (c *Core) GetUserInfo(ctx context.Context, token string) (*user.UserDTO, error) {
	session, err := c.SessionsRepo.Get(ctx, token)
	if err != nil {
		return nil, err
	}

	if time.Now().After(session.ExpiresAt) {
		return nil, custom_errors.ErrJWTTimeExpired
	}

	return c.UsersRepo.GetByID(ctx, session.UserID)
}
