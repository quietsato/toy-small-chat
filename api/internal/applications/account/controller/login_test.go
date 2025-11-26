package controller_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/quietsato/toy-small-chat/api/internal/applications/account/controller"
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

func TestLoginController_Login(t *testing.T) {
	t.Parallel()

	// 事前にハッシュ化されたパスワードを準備
	rawPassword, _ := domain.NewRawPassword([]byte("testpass123"))
	hashedPassword, _ := domain.NewHashedPassword(rawPassword)

	t.Run("ログイン成功", func(t *testing.T) {
		t.Parallel()

		mockQP := &mockAccountQueryProcessor{
			getLoginCredentialFunc: func(ctx context.Context, inp queryprocessor.GetLoginCredentialInput) (queryprocessor.GetLoginCredentialOutput, error) {
				return queryprocessor.GetLoginCredentialOutput{
					AccountID:    uuid.New(),
					PasswordHash: hashedPassword.Bytes(),
				}, nil
			},
		}

		ctrl := controller.NewLoginController(mockQP)

		out, err := ctrl.Login(t.Context(), controller.LoginInput{
			UserName: "testuser",
			Password: "testpass123",
		})

		require.NoError(t, err)
		require.Equal(t, "testuser", out.UserName)
		require.NotEmpty(t, out.Token)
	})

	t.Run("無効なユーザー名でエラーを返す", func(t *testing.T) {
		t.Parallel()

		mockQP := &mockAccountQueryProcessor{}

		ctrl := controller.NewLoginController(mockQP)

		_, err := ctrl.Login(t.Context(), controller.LoginInput{
			UserName: "", // invalid
			Password: "testpass123",
		})

		require.Error(t, err)
	})

	t.Run("無効なパスワードでエラーを返す", func(t *testing.T) {
		t.Parallel()

		mockQP := &mockAccountQueryProcessor{}

		ctrl := controller.NewLoginController(mockQP)

		_, err := ctrl.Login(t.Context(), controller.LoginInput{
			UserName: "testuser",
			Password: "short", // too short
		})

		require.Error(t, err)
	})

	t.Run("アカウントが存在しない場合にエラーを返す", func(t *testing.T) {
		t.Parallel()

		mockQP := &mockAccountQueryProcessor{
			getLoginCredentialFunc: func(ctx context.Context, inp queryprocessor.GetLoginCredentialInput) (queryprocessor.GetLoginCredentialOutput, error) {
				return queryprocessor.GetLoginCredentialOutput{}, queryprocessor.ErrAccountNotFound
			},
		}

		ctrl := controller.NewLoginController(mockQP)

		_, err := ctrl.Login(t.Context(), controller.LoginInput{
			UserName: "nonexistent",
			Password: "testpass123",
		})

		require.Error(t, err)
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

		ctrl := controller.NewLoginController(mockQP)

		_, err := ctrl.Login(t.Context(), controller.LoginInput{
			UserName: "testuser",
			Password: "wrongpass1",
		})

		require.Error(t, err)
	})
}

func TestNewLoginController(t *testing.T) {
	t.Parallel()

	t.Run("正しく初期化される", func(t *testing.T) {
		t.Parallel()
		mockQP := &mockAccountQueryProcessor{}

		ctrl := controller.NewLoginController(mockQP)

		require.NotNil(t, ctrl)
	})
}
