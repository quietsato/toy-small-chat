package controller_test

import (
	"context"
	"errors"
	"testing"

	"github.com/quietsato/toy-small-chat/api/internal/applications/room/controller"
	"github.com/quietsato/toy-small-chat/api/internal/applications/room/usecase/queryprocessor"
	"github.com/stretchr/testify/require"
)

// Mock implementations
type mockRoomQueryProcessor struct {
	getRoomsFunc func(ctx context.Context, inp queryprocessor.GetRoomsInput) (queryprocessor.GetRoomsOutput, error)
}

func (m *mockRoomQueryProcessor) GetRooms(ctx context.Context, inp queryprocessor.GetRoomsInput) (queryprocessor.GetRoomsOutput, error) {
	if m.getRoomsFunc != nil {
		return m.getRoomsFunc(ctx, inp)
	}
	return queryprocessor.GetRoomsOutput{}, nil
}

func TestGetRoomsController_GetRooms(t *testing.T) {
	t.Parallel()

	t.Run("ルーム一覧取得成功", func(t *testing.T) {
		t.Parallel()

		mockQP := &mockRoomQueryProcessor{
			getRoomsFunc: func(ctx context.Context, inp queryprocessor.GetRoomsInput) (queryprocessor.GetRoomsOutput, error) {
				return queryprocessor.GetRoomsOutput{
					Rooms: []queryprocessor.RoomDTO{
						{
							ID:        "room-1",
							Name:      "General",
							CreatedBy: "user-1",
							CreatedAt: "2024-01-01T00:00:00Z",
							UpdatedAt: "2024-01-01T00:00:00Z",
						},
						{
							ID:        "room-2",
							Name:      "Random",
							CreatedBy: "user-2",
							CreatedAt: "2024-01-02T00:00:00Z",
							UpdatedAt: "2024-01-02T00:00:00Z",
						},
					},
				}, nil
			},
		}

		ctrl := controller.NewGetRoomsController(mockQP)

		out, err := ctrl.GetRooms(t.Context(), controller.GetRoomsInput{})

		require.NoError(t, err)
		require.Len(t, out.Rooms, 2)
		require.Equal(t, "room-1", out.Rooms[0].ID)
		require.Equal(t, "General", out.Rooms[0].Name)
	})

	t.Run("空のルーム一覧", func(t *testing.T) {
		t.Parallel()

		mockQP := &mockRoomQueryProcessor{
			getRoomsFunc: func(ctx context.Context, inp queryprocessor.GetRoomsInput) (queryprocessor.GetRoomsOutput, error) {
				return queryprocessor.GetRoomsOutput{
					Rooms: []queryprocessor.RoomDTO{},
				}, nil
			},
		}

		ctrl := controller.NewGetRoomsController(mockQP)

		out, err := ctrl.GetRooms(t.Context(), controller.GetRoomsInput{})

		require.NoError(t, err)
		require.Empty(t, out.Rooms)
	})

	t.Run("クエリプロセッサエラー時にエラーを返す", func(t *testing.T) {
		t.Parallel()

		mockQP := &mockRoomQueryProcessor{
			getRoomsFunc: func(ctx context.Context, inp queryprocessor.GetRoomsInput) (queryprocessor.GetRoomsOutput, error) {
				return queryprocessor.GetRoomsOutput{}, errors.New("db error")
			},
		}

		ctrl := controller.NewGetRoomsController(mockQP)

		_, err := ctrl.GetRooms(t.Context(), controller.GetRoomsInput{})

		require.Error(t, err)
	})
}

func TestNewGetRoomsController(t *testing.T) {
	t.Parallel()

	t.Run("正しく初期化される", func(t *testing.T) {
		t.Parallel()
		mockQP := &mockRoomQueryProcessor{}

		ctrl := controller.NewGetRoomsController(mockQP)

		require.NotNil(t, ctrl)
	})
}
