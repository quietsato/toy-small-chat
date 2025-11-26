package controller_test

import (
	"context"
	"errors"
	"testing"

	"github.com/quietsato/toy-small-chat/api/internal/applications/account/controller"
	"github.com/quietsato/toy-small-chat/api/internal/applications/account/usecase/repository"
	"github.com/quietsato/toy-small-chat/api/internal/domain"
	"github.com/stretchr/testify/require"
)

// Mock implementations
type mockAccountRepository struct {
	createAccountFunc func(ctx context.Context, inp repository.CreateAccountInput) (repository.CreateAccountOutput, error)
}

func (m *mockAccountRepository) CreateAccount(ctx context.Context, inp repository.CreateAccountInput) (repository.CreateAccountOutput, error) {
	if m.createAccountFunc != nil {
		return m.createAccountFunc(ctx, inp)
	}
	return repository.CreateAccountOutput{}, nil
}

type mockAuthService struct {
	generateTokenFunc func(id string) string
}

func (m *mockAuthService) GenerateToken(id string) string {
	if m.generateTokenFunc != nil {
		return m.generateTokenFunc(id)
	}
	return "mock-token"
}

func TestCreateAccountController_CreateAccount(t *testing.T) {
	t.Parallel()

	t.Run("アカウント作成成功", func(t *testing.T) {
		t.Parallel()

		mockRepo := &mockAccountRepository{
			createAccountFunc: func(ctx context.Context, inp repository.CreateAccountInput) (repository.CreateAccountOutput, error) {
				return repository.CreateAccountOutput{
					AccountID: "test-account-id",
				}, nil
			},
		}

		mockAuth := &mockAuthService{
			generateTokenFunc: func(id string) string {
				return "generated-token"
			},
		}

		ctrl := controller.NewCreateAccountController(mockRepo, mockAuth)

		out, err := ctrl.CreateAccount(t.Context(), controller.CreateAccountInput{
			UserName: "testuser",
			Password: "testpass123",
		})

		require.NoError(t, err)
		require.Equal(t, "testuser", out.UserName)
		require.Equal(t, "generated-token", out.Token)
	})

	t.Run("無効なユーザー名でエラーを返す", func(t *testing.T) {
		t.Parallel()

		mockRepo := &mockAccountRepository{}
		mockAuth := &mockAuthService{}

		ctrl := controller.NewCreateAccountController(mockRepo, mockAuth)

		_, err := ctrl.CreateAccount(t.Context(), controller.CreateAccountInput{
			UserName: "", // invalid
			Password: "testpass123",
		})

		require.Error(t, err)
		if !errors.Is(err, domain.ErrInvalidUserName) {
			require.NotEmpty(t, err.Error())
		}
	})

	t.Run("無効なパスワードでエラーを返す", func(t *testing.T) {
		t.Parallel()

		mockRepo := &mockAccountRepository{}
		mockAuth := &mockAuthService{}

		ctrl := controller.NewCreateAccountController(mockRepo, mockAuth)

		_, err := ctrl.CreateAccount(t.Context(), controller.CreateAccountInput{
			UserName: "testuser",
			Password: "short", // too short
		})

		require.Error(t, err)
	})

	t.Run("リポジトリエラー時にエラーを返す", func(t *testing.T) {
		t.Parallel()

		mockRepo := &mockAccountRepository{
			createAccountFunc: func(ctx context.Context, inp repository.CreateAccountInput) (repository.CreateAccountOutput, error) {
				return repository.CreateAccountOutput{}, errors.New("db error")
			},
		}

		mockAuth := &mockAuthService{}

		ctrl := controller.NewCreateAccountController(mockRepo, mockAuth)

		_, err := ctrl.CreateAccount(t.Context(), controller.CreateAccountInput{
			UserName: "testuser",
			Password: "testpass123",
		})

		require.Error(t, err)
	})
}

func TestNewCreateAccountController(t *testing.T) {
	t.Parallel()

	t.Run("正しく初期化される", func(t *testing.T) {
		t.Parallel()
		mockRepo := &mockAccountRepository{}
		mockAuth := &mockAuthService{}

		ctrl := controller.NewCreateAccountController(mockRepo, mockAuth)

		require.NotNil(t, ctrl)
	})
}
