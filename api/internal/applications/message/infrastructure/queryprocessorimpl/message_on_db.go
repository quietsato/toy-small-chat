package queryprocessorimpl

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/quietsato/toy-small-chat/api/internal/applications/message/usecase/queryprocessor"
	"github.com/quietsato/toy-small-chat/api/internal/db"
)

type MessageQueryProcessorOnDB struct {
	queries *db.Queries
}

func NewMessageQueryProcessorOnDB(pool *pgxpool.Pool) *MessageQueryProcessorOnDB {
	return &MessageQueryProcessorOnDB{
		queries: db.New(pool),
	}
}

func (q *MessageQueryProcessorOnDB) GetMessages(roomID string) ([]queryprocessor.Message, error) {
	ctx := context.Background()

	dbMessages, err := q.queries.GetMessagesByRoomID(ctx, uuid.MustParse(roomID))
	if err != nil {
		return nil, err
	}

	messages := make([]queryprocessor.Message, 0, len(dbMessages))
	for _, dbMsg := range dbMessages {
		var createdAtStr string
		if dbMsg.CreatedAt.Valid {
			createdAtStr = dbMsg.CreatedAt.Time.Format("2006-01-02T15:04:05Z07:00")
		}

		messages = append(messages, queryprocessor.Message{
			ID:        dbMsg.MessageID.String(),
			Author:    dbMsg.AuthorName,
			Content:   dbMsg.Content,
			CreatedAt: createdAtStr,
		})
	}

	return messages, nil
}

var _ queryprocessor.MessageQueryProcessor = (*MessageQueryProcessorOnDB)(nil)
