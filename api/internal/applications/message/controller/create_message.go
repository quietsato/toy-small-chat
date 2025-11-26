package controller

import (
	"context"

	"github.com/quietsato/toy-small-chat/api/internal/applications/message/usecase/repository"
)

type CreateMessageController struct {
	repo repository.MessageRepository
}

func NewCreateMessageController(repo repository.MessageRepository) *CreateMessageController {
	return &CreateMessageController{repo}
}

func (c *CreateMessageController) CreateMessage(ctx context.Context, inp CreateMessageInput) error {
	err := c.repo.CreateMessage(ctx, repository.CreateMessageInput{
		AuthorID: inp.AuthorID,
		RoomID:   inp.RoomID,
		Content:  inp.Content,
	})
	if err != nil {
		return err
	}

	return nil
}

type CreateMessageInput struct {
	RoomID   string `json:"roomID"`
	Content  string `json:"content"`
	AuthorID string `json:"-"`
}

type CreateMessageOutput struct{}
