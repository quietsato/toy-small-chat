package usecase

import (
	"context"
	"errors"
	"fmt"

	"github.com/quietsato/toy-small-chat/api/internal/applications/account/usecase/queryprocessor"
	"github.com/quietsato/toy-small-chat/api/internal/applications/account/usecase/service"
	"github.com/quietsato/toy-small-chat/api/internal/domain"
)

type LoginInput struct {
	UserName domain.UserName
	Password domain.RawPassword
}
type LoginOutput struct {
	Token string // JWT はドメインモデルとは言い難いので直接 string として扱う
}

var (
	ErrAccountNotFound    = errors.New("account not found")
	ErrPasswordIsNotMatch = errors.New("password not match")
)

func NewLoginUsecase(q queryprocessor.AccountQueryProcessor, authService service.AuthService) *LoginUsecase {
	return &LoginUsecase{q, authService}
}

type LoginUsecase struct {
	q    queryprocessor.AccountQueryProcessor
	auth service.AuthService
}

func (u *LoginUsecase) Execute(ctx context.Context, inp LoginInput) (LoginOutput, error) {
	res, err := u.q.GetLoginCredential(ctx, queryprocessor.GetLoginCredentialInput{
		UserName: inp.UserName.String(),
	})
	if errors.Is(err, queryprocessor.ErrAccountNotFound) {
		return LoginOutput{}, ErrAccountNotFound
	}
	if err != nil {
		return LoginOutput{}, fmt.Errorf("failed to login: %w", err)
	}

	pass := domain.NewHashedPasswordFromHash([]byte(res.PasswordHash))
	matched, err := pass.Match(inp.Password)
	if err != nil {
		return LoginOutput{}, err
	}
	if !matched {
		return LoginOutput{}, ErrPasswordIsNotMatch
	}

	token := u.auth.GenerateToken(res.AccountID.String())

	return LoginOutput{
		Token: token,
	}, nil
}
