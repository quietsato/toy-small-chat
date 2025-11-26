package usecase_test

import (
	"context"
	"errors"
	"testing"

	"github.com/quietsato/toy-small-chat/api/internal/applications/room/usecase"
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

func TestCreateRoomUsecase_Execute(t *testing.T) {
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

		uc := usecase.NewCreateRoomUsecase(mockRepo)

		_, err := uc.Execute(t.Context(), usecase.CreateRoomInput{
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

		uc := usecase.NewCreateRoomUsecase(mockRepo)

		_, err := uc.Execute(t.Context(), usecase.CreateRoomInput{
			Name:      "test-room",
			CreatedBy: "user-123",
		})

		require.Error(t, err)
	})
}

func TestNewCreateRoomUsecase(t *testing.T) {
	t.Parallel()

	t.Run("正しく初期化される", func(t *testing.T) {
		t.Parallel()
		mockRepo := &mockRoomRepository{}

		uc := usecase.NewCreateRoomUsecase(mockRepo)

		require.NotNil(t, uc)
	})
}
