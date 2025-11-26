package queryprocessorimpl

import (
	"context"

	"github.com/quietsato/toy-small-chat/api/internal/applications/room/usecase/queryprocessor"
)

type MockRoomQueryProcessor struct{}

func (m *MockRoomQueryProcessor) GetRooms(ctx context.Context, inp queryprocessor.GetRoomsInput) (queryprocessor.GetRoomsOutput, error) {
	return queryprocessor.GetRoomsOutput{
		Rooms: []queryprocessor.RoomDTO{},
	}, nil
}

var _ queryprocessor.RoomQueryProcessor = new(MockRoomQueryProcessor)
