package controller_test

import (
	"context"
	"errors"
	"testing"

	"github.com/quietsato/toy-small-chat/api/internal/applications/message/controller"
	"github.com/quietsato/toy-small-chat/api/internal/applications/message/usecase/repository"
	"github.com/stretchr/testify/require"
)

// Mock implementations
type mockMessageRepository struct {
	createMessageFunc func(ctx context.Context, inp repository.CreateMessageInput) error
	getFunc           func() error
}

func (m *mockMessageRepository) CreateMessage(ctx context.Context, inp repository.CreateMessageInput) error {
	if m.createMessageFunc != nil {
		return m.createMessageFunc(ctx, inp)
	}
	return nil
}

func (m *mockMessageRepository) Get() error {
	if m.getFunc != nil {
		return m.getFunc()
	}
	return nil
}

func TestCreateMessageController_CreateMessage(t *testing.T) {
	t.Parallel()

	t.Run("メッセージ作成成功", func(t *testing.T) {
		t.Parallel()

		mockRepo := &mockMessageRepository{
			createMessageFunc: func(ctx context.Context, inp repository.CreateMessageInput) error {
				require.Equal(t, "author-123", inp.AuthorID)
				require.Equal(t, "room-456", inp.RoomID)
				require.Equal(t, "Hello, World!", inp.Content)
				return nil
			},
		}

		ctrl := controller.NewCreateMessageController(mockRepo)

		err := ctrl.CreateMessage(t.Context(), controller.CreateMessageInput{
			AuthorID: "author-123",
			RoomID:   "room-456",
			Content:  "Hello, World!",
		})

		require.NoError(t, err)
	})

	t.Run("リポジトリエラー時にエラーを返す", func(t *testing.T) {
		t.Parallel()

		mockRepo := &mockMessageRepository{
			createMessageFunc: func(ctx context.Context, inp repository.CreateMessageInput) error {
				return errors.New("db error")
			},
		}

		ctrl := controller.NewCreateMessageController(mockRepo)

		err := ctrl.CreateMessage(t.Context(), controller.CreateMessageInput{
			AuthorID: "author-123",
			RoomID:   "room-456",
			Content:  "Hello, World!",
		})

		require.Error(t, err)
	})
}

func TestNewCreateMessageController(t *testing.T) {
	t.Parallel()

	t.Run("正しく初期化される", func(t *testing.T) {
		t.Parallel()
		mockRepo := &mockMessageRepository{}

		ctrl := controller.NewCreateMessageController(mockRepo)

		require.NotNil(t, ctrl)
	})
}
