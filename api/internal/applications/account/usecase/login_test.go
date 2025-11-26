package usecase_test

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/quietsato/toy-small-chat/api/internal/applications/account/usecase"
	"github.com/quietsato/toy-small-chat/api/internal/applications/account/usecase/queryprocessor"
	"github.com/quietsato/toy-small-chat/api/internal/domain"
	"github.com/stretchr/testify/require"
)

// Mock implementations
type mockAccountQueryProcessor struct {
	getLoginCredentialFunc func(ctx context.Context, inp queryprocessor.GetLoginCredentialInput) (queryprocessor.GetLoginCredentialOutput, error)
}

func (m *mockAccountQueryProcessor) GetLoginCredential(ctx context.Context, inp queryprocessor.GetLoginCredentialInput) (queryprocessor.GetLoginCredentialOutput, error) {
	if m.getLoginCredentialFunc != nil {
		return m.getLoginCredentialFunc(ctx, inp)
	}
	return queryprocessor.GetLoginCredentialOutput{}, nil
}

func TestLoginUsecase_Execute(t *testing.T) {
	t.Parallel()

	// 事前にハッシュ化されたパスワードを準備
	rawPassword, _ := domain.NewRawPassword([]byte("testpass123"))
	hashedPassword, _ := domain.NewHashedPassword(rawPassword)

	t.Run("ログイン成功", func(t *testing.T) {
		t.Parallel()

		accountID := uuid.New()

		mockQP := &mockAccountQueryProcessor{
			getLoginCredentialFunc: func(ctx context.Context, inp queryprocessor.GetLoginCredentialInput) (queryprocessor.GetLoginCredentialOutput, error) {
				require.Equal(t, "testuser", inp.UserName)
				return queryprocessor.GetLoginCredentialOutput{
					AccountID:    accountID,
					PasswordHash: hashedPassword.Bytes(),
				}, nil
			},
		}

		mockAuth := &mockAuthService{
			generateTokenFunc: func(id string) string {
				require.Equal(t, accountID.String(), id)
				return "generated-token"
			},
		}

		uc := usecase.NewLoginUsecase(mockQP, mockAuth)

		userName, _ := domain.NewUserName("testuser")
		password, _ := domain.NewRawPassword([]byte("testpass123"))

		out, err := uc.Execute(t.Context(), usecase.LoginInput{
			UserName: userName,
			Password: password,
		})

		require.NoError(t, err)
		require.Equal(t, "generated-token", out.Token)
	})

	t.Run("アカウントが存在しない場合にエラーを返す", func(t *testing.T) {
		t.Parallel()

		mockQP := &mockAccountQueryProcessor{
			getLoginCredentialFunc: func(ctx context.Context, inp queryprocessor.GetLoginCredentialInput) (queryprocessor.GetLoginCredentialOutput, error) {
				return queryprocessor.GetLoginCredentialOutput{}, queryprocessor.ErrAccountNotFound
			},
		}

		mockAuth := &mockAuthService{}

		uc := usecase.NewLoginUsecase(mockQP, mockAuth)

		userName, _ := domain.NewUserName("nonexistent")
		password, _ := domain.NewRawPassword([]byte("testpass123"))

		_, err := uc.Execute(t.Context(), usecase.LoginInput{
			UserName: userName,
			Password: password,
		})

		require.ErrorIs(t, err, usecase.ErrAccountNotFound)
	})

	t.Run("パスワードが一致しない場合にエラーを返す", func(t *testing.T) {
		t.Parallel()

		mockQP := &mockAccountQueryProcessor{
			getLoginCredentialFunc: func(ctx context.Context, inp queryprocessor.GetLoginCredentialInput) (queryprocessor.GetLoginCredentialOutput, error) {
				return queryprocessor.GetLoginCredentialOutput{
					AccountID:    uuid.New(),
					PasswordHash: hashedPassword.Bytes(),
				}, nil
			},
		}

		mockAuth := &mockAuthService{}

		uc := usecase.NewLoginUsecase(mockQP, mockAuth)

		userName, _ := domain.NewUserName("testuser")
		wrongPassword, _ := domain.NewRawPassword([]byte("wrongpass1"))

		_, err := uc.Execute(t.Context(), usecase.LoginInput{
			UserName: userName,
			Password: wrongPassword,
		})

		require.ErrorIs(t, err, usecase.ErrPasswordIsNotMatch)
	})

	t.Run("クエリプロセッサのエラー時にエラーを返す", func(t *testing.T) {
		t.Parallel()

		mockQP := &mockAccountQueryProcessor{
			getLoginCredentialFunc: func(ctx context.Context, inp queryprocessor.GetLoginCredentialInput) (queryprocessor.GetLoginCredentialOutput, error) {
				return queryprocessor.GetLoginCredentialOutput{}, errors.New("db error")
			},
		}

		mockAuth := &mockAuthService{}

		uc := usecase.NewLoginUsecase(mockQP, mockAuth)

		userName, _ := domain.NewUserName("testuser")
		password, _ := domain.NewRawPassword([]byte("testpass123"))

		_, err := uc.Execute(t.Context(), usecase.LoginInput{
			UserName: userName,
			Password: password,
		})

		require.Error(t, err)
	})
}

func TestNewLoginUsecase(t *testing.T) {
	t.Parallel()

	t.Run("正しく初期化される", func(t *testing.T) {
		t.Parallel()
		mockQP := &mockAccountQueryProcessor{}
		mockAuth := &mockAuthService{}

		uc := usecase.NewLoginUsecase(mockQP, mockAuth)

		require.NotNil(t, uc)
	})
}
