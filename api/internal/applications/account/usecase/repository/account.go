package repository

import (
	"context"
	"errors"
)

type CreateAccountInput struct {
	UserName     string
	PasswordHash []byte
}
type CreateAccountOutput struct {
	AccountID string
}

var (
	ErrUserNameAlreadyRegistered = errors.New("username already registered")
)

type AccountRepository interface {
	CreateAccount(ctx context.Context, inp CreateAccountInput) (CreateAccountOutput, error)
}
