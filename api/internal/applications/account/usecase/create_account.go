package usecase

import (
	"context"
	"fmt"

	"github.com/quietsato/toy-small-chat/api/internal/applications/account/usecase/repository"
	"github.com/quietsato/toy-small-chat/api/internal/applications/account/usecase/service"
	"github.com/quietsato/toy-small-chat/api/internal/domain"
)

type CreateAccountUsecase struct {
	r    repository.AccountRepository
	auth service.AuthService
}

type CreateAccountInput struct {
	UserName domain.UserName
	Password domain.RawPassword
}

type CreateAccountOutput struct {
	Token string
}

func NewCreateAccountUsecase(r repository.AccountRepository, auth service.AuthService) *CreateAccountUsecase {
	return &CreateAccountUsecase{r, auth}
}

func (u *CreateAccountUsecase) Execute(ctx context.Context, inp CreateAccountInput) (CreateAccountOutput, error) {
	hashed, err := domain.NewHashedPassword(inp.Password)
	if err != nil {
		return CreateAccountOutput{}, err
	}

	res, err := u.r.CreateAccount(ctx, repository.CreateAccountInput{
		UserName:     inp.UserName.String(),
		PasswordHash: hashed.Bytes(),
	})
	if err != nil {
		return CreateAccountOutput{}, fmt.Errorf("failed to create account: %w", err)
	}

	token := u.auth.GenerateToken(res.AccountID)
	return CreateAccountOutput{
		Token: token,
	}, nil
}
