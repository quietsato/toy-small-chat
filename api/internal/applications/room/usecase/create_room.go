package usecase

import (
	"context"
	"fmt"

	"github.com/quietsato/toy-small-chat/api/internal/applications/room/usecase/repository"
)

type CreateRoomUsecase struct {
	repo repository.RoomRepository
}

type CreateRoomInput struct {
	Name      string
	CreatedBy string
}

type CreateRoomOutput struct{}

func NewCreateRoomUsecase(repo repository.RoomRepository) *CreateRoomUsecase {
	return &CreateRoomUsecase{repo}
}

func (u *CreateRoomUsecase) Execute(ctx context.Context, inp CreateRoomInput) (CreateRoomOutput, error) {
	_, err := u.repo.CreateRoom(ctx, repository.CreateRoomInput{
		Name:      inp.Name,
		CreatedBy: inp.CreatedBy,
	})
	if err != nil {
		return CreateRoomOutput{}, fmt.Errorf("failed to create room: %w", err)
	}

	return CreateRoomOutput{}, nil
}
