package controller

import (
	"context"
	"fmt"

	"github.com/quietsato/toy-small-chat/api/internal/applications/account/infrastructure/serviceimpl"
	"github.com/quietsato/toy-small-chat/api/internal/applications/account/usecase"
	"github.com/quietsato/toy-small-chat/api/internal/applications/account/usecase/queryprocessor"
	"github.com/quietsato/toy-small-chat/api/internal/domain"
)

type LoginInput struct {
	UserName string `json:"username"`
	Password string `json:"password"`
}
type LoginOutput struct {
	UserName string `json:"username"`
	Token    string `json:"token"`
}

type LoginController struct {
	query queryprocessor.AccountQueryProcessor
}

func NewLoginController(query queryprocessor.AccountQueryProcessor) *LoginController {
	return &LoginController{query}
}

func (c *LoginController) Login(ctx context.Context, inp LoginInput) (LoginOutput, error) {
	auth := serviceimpl.NewAuthService([]byte("secretKey"))

	userName, err := domain.NewUserName(inp.UserName)
	if err != nil {
		return LoginOutput{}, fmt.Errorf("bad username: %w", err)
	}
	password, err := domain.NewRawPassword([]byte(inp.Password))
	if err != nil {
		return LoginOutput{}, fmt.Errorf("bad password: %w", err)
	}

	uc := usecase.NewLoginUsecase(c.query, auth)
	res, err := uc.Execute(ctx, usecase.LoginInput{
		UserName: userName,
		Password: password,
	})
	if err != nil {
		return LoginOutput{}, fmt.Errorf("failed to login: %w", err)
	}

	return LoginOutput{
		UserName: inp.UserName,
		Token:    res.Token,
	}, nil
}
