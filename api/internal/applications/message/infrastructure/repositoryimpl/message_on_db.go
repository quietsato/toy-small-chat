package repositoryimpl

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/quietsato/toy-small-chat/api/internal/applications/message/usecase/repository"
	"github.com/quietsato/toy-small-chat/api/internal/db"
)

type MessageRepositoryOnDB struct {
	pool *pgxpool.Pool
}

func NewMessageRepositoryOnDB(pool *pgxpool.Pool) *MessageRepositoryOnDB {
	return &MessageRepositoryOnDB{pool}
}

func (r *MessageRepositoryOnDB) CreateMessage(ctx context.Context, inp repository.CreateMessageInput) error {

	// Begin Transaction
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() {
		if err := tx.Rollback(ctx); err != nil {
			slog.ErrorContext(ctx, "failed to rollback", slog.Any("err", err))
		}
	}()

	// Exec
	queries := db.New(r.pool).WithTx(tx)

	if err := queries.CreateMessage(ctx, db.CreateMessageParams{
		AuthorID: uuid.MustParse(inp.AuthorID),
		Content:  inp.Content,
		RoomID:   uuid.MustParse(inp.RoomID),
	}); err != nil {
		return fmt.Errorf("failed to create message: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit: %w", err)
	}

	return nil
}

func (r *MessageRepositoryOnDB) Get() error {
	panic("unimplemented")
}

var _ repository.MessageRepository = new(MessageRepositoryOnDB)
