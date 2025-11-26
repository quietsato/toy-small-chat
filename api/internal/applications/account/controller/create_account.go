package controller

import (
	"context"
	"fmt"

	"github.com/quietsato/toy-small-chat/api/internal/applications/account/usecase"
	"github.com/quietsato/toy-small-chat/api/internal/applications/account/usecase/repository"
	"github.com/quietsato/toy-small-chat/api/internal/applications/account/usecase/service"
	"github.com/quietsato/toy-small-chat/api/internal/domain"
)

type CreateAccountInput struct {
	UserName string `json:"username"`
	Password string `json:"password"`
}
type CreateAccountOutput struct {
	UserName string `json:"username"`
	Token    string `json:"token"`
}

type CreateAccountController struct {
	repo repository.AccountRepository
	auth service.AuthService
}

func NewCreateAccountController(repo repository.AccountRepository, auth service.AuthService) *CreateAccountController {
	return &CreateAccountController{repo, auth}
}

func (c *CreateAccountController) CreateAccount(ctx context.Context, inp CreateAccountInput) (CreateAccountOutput, error) {
	userName, err := domain.NewUserName(inp.UserName)
	if err != nil {
		return CreateAccountOutput{}, fmt.Errorf("bad username: %w", err)
	}
	password, err := domain.NewRawPassword([]byte(inp.Password))
	if err != nil {
		return CreateAccountOutput{}, fmt.Errorf("bad password: %w", err)
	}

	uc := usecase.NewCreateAccountUsecase(c.repo, c.auth)
	res, err := uc.Execute(ctx, usecase.CreateAccountInput{
		UserName: userName,
		Password: password,
	})
	if err != nil {
		return CreateAccountOutput{}, fmt.Errorf("failed to create account: %w", err)
	}

	return CreateAccountOutput{
		UserName: inp.UserName,
		Token:    res.Token,
	}, nil
}
