package controller_test

import (
	"errors"
	"testing"

	"github.com/quietsato/toy-small-chat/api/internal/applications/message/controller"
	"github.com/quietsato/toy-small-chat/api/internal/applications/message/usecase/queryprocessor"
	"github.com/stretchr/testify/require"
)

// Mock implementations
type mockMessageQueryProcessor struct {
	getMessagesFunc func(roomID string) ([]queryprocessor.Message, error)
}

func (m *mockMessageQueryProcessor) GetMessages(roomID string) ([]queryprocessor.Message, error) {
	if m.getMessagesFunc != nil {
		return m.getMessagesFunc(roomID)
	}
	return nil, nil
}

func TestGetMessagesController_GetMessages(t *testing.T) {
	t.Parallel()

	t.Run("メッセージ取得成功", func(t *testing.T) {
		t.Parallel()

		mockQP := &mockMessageQueryProcessor{
			getMessagesFunc: func(roomID string) ([]queryprocessor.Message, error) {
				require.Equal(t, "room-123", roomID)
				return []queryprocessor.Message{
					{
						ID:        "msg-1",
						Author:    "user-1",
						Content:   "Hello",
						CreatedAt: "2024-01-01T00:00:00Z",
					},
					{
						ID:        "msg-2",
						Author:    "user-2",
						Content:   "World",
						CreatedAt: "2024-01-01T00:01:00Z",
					},
				}, nil
			},
		}

		ctrl := controller.NewGetMessagesController(mockQP)

		out, err := ctrl.GetMessages(controller.GetMessagesInput{
			RoomID: "room-123",
		})

		require.NoError(t, err)
		require.Len(t, out.Messages, 2)
		require.Equal(t, "msg-1", out.Messages[0].ID)
		require.Equal(t, "Hello", out.Messages[0].Content)
	})

	t.Run("空のメッセージリスト", func(t *testing.T) {
		t.Parallel()

		mockQP := &mockMessageQueryProcessor{
			getMessagesFunc: func(roomID string) ([]queryprocessor.Message, error) {
				return []queryprocessor.Message{}, nil
			},
		}

		ctrl := controller.NewGetMessagesController(mockQP)

		out, err := ctrl.GetMessages(controller.GetMessagesInput{
			RoomID: "room-123",
		})

		require.NoError(t, err)
		require.Empty(t, out.Messages)
	})

	t.Run("クエリプロセッサエラー時にエラーを返す", func(t *testing.T) {
		t.Parallel()

		mockQP := &mockMessageQueryProcessor{
			getMessagesFunc: func(roomID string) ([]queryprocessor.Message, error) {
				return nil, errors.New("db error")
			},
		}

		ctrl := controller.NewGetMessagesController(mockQP)

		_, err := ctrl.GetMessages(controller.GetMessagesInput{
			RoomID: "room-123",
		})

		require.Error(t, err)
	})
}

func TestNewGetMessagesController(t *testing.T) {
	t.Parallel()

	t.Run("正しく初期化される", func(t *testing.T) {
		t.Parallel()
		mockQP := &mockMessageQueryProcessor{}

		ctrl := controller.NewGetMessagesController(mockQP)

		require.NotNil(t, ctrl)
	})
}
