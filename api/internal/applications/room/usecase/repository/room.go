package repository

import "context"

type CreateRoomInput struct {
	Name      string
	CreatedBy string
}
type CreateRoomOutput struct{}

type RoomRepository interface {
	CreateRoom(ctx context.Context, inp CreateRoomInput) (CreateRoomOutput, error)
}
