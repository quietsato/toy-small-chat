package controller

import (
	"context"
	"fmt"

	"github.com/quietsato/toy-small-chat/api/internal/applications/room/usecase/queryprocessor"
)

type GetRoomsInput struct{}

type GetRoomsOutput struct {
	Rooms []Room `json:"rooms"`
}

type Room struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	CreatedBy string `json:"createdBy"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}

type GetRoomsController struct {
	query queryprocessor.RoomQueryProcessor
}

func NewGetRoomsController(query queryprocessor.RoomQueryProcessor) *GetRoomsController {
	return &GetRoomsController{query}
}

func (c *GetRoomsController) GetRooms(ctx context.Context, inp GetRoomsInput) (GetRoomsOutput, error) {
	res, err := c.query.GetRooms(ctx, queryprocessor.GetRoomsInput{})
	if err != nil {
		return GetRoomsOutput{}, fmt.Errorf("failed to get rooms: %w", err)
	}

	rooms := make([]Room, 0, len(res.Rooms))
	for _, dto := range res.Rooms {
		rooms = append(rooms, Room{
			ID:        dto.ID,
			Name:      dto.Name,
			CreatedBy: dto.CreatedBy,
			CreatedAt: dto.CreatedAt,
			UpdatedAt: dto.UpdatedAt,
		})
	}

	return GetRoomsOutput{
		Rooms: rooms,
	}, nil
}
