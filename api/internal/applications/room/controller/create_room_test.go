package controller_test

import (
	"context"
	"errors"
	"testing"

	"github.com/quietsato/toy-small-chat/api/internal/applications/room/controller"
	"github.com/quietsato/toy-small-chat/api/internal/applications/room/usecase/repository"
	"github.com/stretchr/testify/require"
)

// Mock implementations
type mockRoomRepository struct {
	createRoomFunc func(ctx context.Context, inp repository.CreateRoomInput) (repository.CreateRoomOutput, error)
}

func (m *mockRoomRepository) CreateRoom(ctx context.Context, inp repository.CreateRoomInput) (repository.CreateRoomOutput, error) {
	if m.createRoomFunc != nil {
		return m.createRoomFunc(ctx, inp)
	}
	return repository.CreateRoomOutput{}, nil
}

func TestCreateRoomController_CreateRoom(t *testing.T) {
	t.Parallel()

	t.Run("ルーム作成成功", func(t *testing.T) {
		t.Parallel()

		mockRepo := &mockRoomRepository{
			createRoomFunc: func(ctx context.Context, inp repository.CreateRoomInput) (repository.CreateRoomOutput, error) {
				require.Equal(t, "test-room", inp.Name)
				require.Equal(t, "user-123", inp.CreatedBy)
				return repository.CreateRoomOutput{}, nil
			},
		}

		ctrl := controller.NewCreateRoomController(mockRepo)

		_, err := ctrl.CreateRoom(t.Context(), controller.CreateRoomInput{
			Name:      "test-room",
			CreatedBy: "user-123",
		})

		require.NoError(t, err)
	})

	t.Run("リポジトリエラー時にエラーを返す", func(t *testing.T) {
		t.Parallel()

		mockRepo := &mockRoomRepository{
			createRoomFunc: func(ctx context.Context, inp repository.CreateRoomInput) (repository.CreateRoomOutput, error) {
				return repository.CreateRoomOutput{}, errors.New("db error")
			},
		}

		ctrl := controller.NewCreateRoomController(mockRepo)

		_, err := ctrl.CreateRoom(t.Context(), controller.CreateRoomInput{
			Name:      "test-room",
			CreatedBy: "user-123",
		})

		require.Error(t, err)
	})
}

func TestNewCreateRoomController(t *testing.T) {
	t.Parallel()

	t.Run("正しく初期化される", func(t *testing.T) {
		t.Parallel()
		mockRepo := &mockRoomRepository{}

		ctrl := controller.NewCreateRoomController(mockRepo)

		require.NotNil(t, ctrl)
	})
}
