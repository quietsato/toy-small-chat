package repositoryimpl

import (
	"context"

	"github.com/google/uuid"
	"github.com/quietsato/toy-small-chat/api/internal/applications/message/usecase/repository"
)

type Message struct {
	ID        string
	Content   string
	Author    string
	CreatedAt string
}

type InMemoryMessageRepository struct {
	msgs *[]Message
}

func NewInMemoryMessageRepository(ctx context.Context, msgs *[]Message) *InMemoryMessageRepository {
	return &InMemoryMessageRepository{msgs}
}

func (m *InMemoryMessageRepository) CreateMessage(ctx context.Context, inp repository.CreateMessageInput) error {
	*m.msgs = append(*m.msgs, Message{
		ID:        uuid.NewString(),
		Content:   inp.Content,
		Author:    "sample-user",
		CreatedAt: "today",
	})

	return nil
}

func (m InMemoryMessageRepository) Get() error {
	panic("unimplemented")
}

var _ repository.MessageRepository = new(InMemoryMessageRepository)
