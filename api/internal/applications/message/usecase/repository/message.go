package repository

import "context"

type CreateMessageInput struct {
	AuthorID string
	Content  string
	RoomID   string
}

type MessageRepository interface {
	CreateMessage(ctx context.Context, inp CreateMessageInput) error
	Get() error
}
