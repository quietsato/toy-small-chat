package queryprocessor

import (
	"context"
	"errors"

	"github.com/google/uuid"
)

type GetLoginCredentialInput struct{ UserName string }
type GetLoginCredentialOutput struct {
	AccountID    uuid.UUID
	PasswordHash []byte
}

var (
	ErrAccountNotFound = errors.New("account not found")
)

type AccountQueryProcessor interface {
	GetLoginCredential(ctx context.Context, inp GetLoginCredentialInput) (GetLoginCredentialOutput, error)
}
