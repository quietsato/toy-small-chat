package usecase_test

import (
	"context"
	"errors"
	"testing"

	"github.com/quietsato/toy-small-chat/api/internal/applications/account/usecase"
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

func TestCreateAccountUsecase_Execute(t *testing.T) {
	t.Parallel()

	t.Run("アカウント作成成功", func(t *testing.T) {
		t.Parallel()

		mockRepo := &mockAccountRepository{
			createAccountFunc: func(ctx context.Context, inp repository.CreateAccountInput) (repository.CreateAccountOutput, error) {
				require.Equal(t, "testuser", inp.UserName)
				require.NotEmpty(t, inp.PasswordHash)
				return repository.CreateAccountOutput{
					AccountID: "test-account-id",
				}, nil
			},
		}

		mockAuth := &mockAuthService{
			generateTokenFunc: func(id string) string {
				require.Equal(t, "test-account-id", id)
				return "generated-token"
			},
		}

		uc := usecase.NewCreateAccountUsecase(mockRepo, mockAuth)

		userName, _ := domain.NewUserName("testuser")
		password, _ := domain.NewRawPassword([]byte("testpass123"))

		out, err := uc.Execute(t.Context(), usecase.CreateAccountInput{
			UserName: userName,
			Password: password,
		})

		require.NoError(t, err)
		require.Equal(t, "generated-token", out.Token)
	})

	t.Run("リポジトリエラー時にエラーを返す", func(t *testing.T) {
		t.Parallel()

		mockRepo := &mockAccountRepository{
			createAccountFunc: func(ctx context.Context, inp repository.CreateAccountInput) (repository.CreateAccountOutput, error) {
				return repository.CreateAccountOutput{}, errors.New("db error")
			},
		}

		mockAuth := &mockAuthService{}

		uc := usecase.NewCreateAccountUsecase(mockRepo, mockAuth)

		userName, _ := domain.NewUserName("testuser")
		password, _ := domain.NewRawPassword([]byte("testpass123"))

		_, err := uc.Execute(t.Context(), usecase.CreateAccountInput{
			UserName: userName,
			Password: password,
		})

		require.Error(t, err)
	})

	t.Run("ユーザー名が既に登録されている場合にエラーを返す", func(t *testing.T) {
		t.Parallel()

		mockRepo := &mockAccountRepository{
			createAccountFunc: func(ctx context.Context, inp repository.CreateAccountInput) (repository.CreateAccountOutput, error) {
				return repository.CreateAccountOutput{}, repository.ErrUserNameAlreadyRegistered
			},
		}

		mockAuth := &mockAuthService{}

		uc := usecase.NewCreateAccountUsecase(mockRepo, mockAuth)

		userName, _ := domain.NewUserName("existinguser")
		password, _ := domain.NewRawPassword([]byte("testpass123"))

		_, err := uc.Execute(t.Context(), usecase.CreateAccountInput{
			UserName: userName,
			Password: password,
		})

		require.Error(t, err)
	})
}

func TestNewCreateAccountUsecase(t *testing.T) {
	t.Parallel()

	t.Run("正しく初期化される", func(t *testing.T) {
		t.Parallel()
		mockRepo := &mockAccountRepository{}
		mockAuth := &mockAuthService{}

		uc := usecase.NewCreateAccountUsecase(mockRepo, mockAuth)

		require.NotNil(t, uc)
	})
}
