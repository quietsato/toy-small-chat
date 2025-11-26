package queryprocessor

import "context"

type GetRoomsInput struct{}
type GetRoomsOutput struct {
	Rooms []RoomDTO
}
type RoomDTO struct {
	ID        string
	Name      string
	CreatedBy string
	CreatedAt string
	UpdatedAt string
}

type RoomQueryProcessor interface {
	GetRooms(ctx context.Context, inp GetRoomsInput) (GetRoomsOutput, error)
}
