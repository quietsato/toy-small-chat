package controller

import (
	"context"
	"fmt"

	"github.com/quietsato/toy-small-chat/api/internal/applications/room/usecase"
	"github.com/quietsato/toy-small-chat/api/internal/applications/room/usecase/repository"
)

type CreateRoomInput struct {
	Name      string `json:"name"`
	CreatedBy string `json:"-"`
}

type CreateRoomOutput struct{}

type CreateRoomController struct {
	repo repository.RoomRepository
}

func NewCreateRoomController(repo repository.RoomRepository) *CreateRoomController {
	return &CreateRoomController{repo}
}

func (c *CreateRoomController) CreateRoom(ctx context.Context, inp CreateRoomInput) (CreateRoomOutput, error) {
	uc := usecase.NewCreateRoomUsecase(c.repo)
	_, err := uc.Execute(ctx, usecase.CreateRoomInput{
		Name:      inp.Name,
		CreatedBy: inp.CreatedBy,
	})
	if err != nil {
		return CreateRoomOutput{}, fmt.Errorf("failed to create room: %w", err)
	}

	return CreateRoomOutput{}, nil
}
