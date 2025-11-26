package repositoryimpl

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/quietsato/toy-small-chat/api/internal/applications/room/usecase/repository"
	"github.com/quietsato/toy-small-chat/api/internal/db"
)

func NewRoomRepositoryOnDB(pool *pgxpool.Pool) *RoomRepositoryOnDB {
	return &RoomRepositoryOnDB{pool}
}

type RoomRepositoryOnDB struct {
	pool *pgxpool.Pool
}

// CreateRoom implements repository.RoomRepository.
func (r *RoomRepositoryOnDB) CreateRoom(ctx context.Context, inp repository.CreateRoomInput) (repository.CreateRoomOutput, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return repository.CreateRoomOutput{}, fmt.Errorf("failed to begin tx: %w", err)
	}
	defer func() {
		if err := tx.Rollback(ctx); err != nil {
			slog.ErrorContext(ctx, "failed to rollback", slog.Any("err", err))
		}
	}()

	queries := db.New(r.pool).WithTx(tx)
	err = queries.CreateRoom(ctx, db.CreateRoomParams{
		Name:      inp.Name,
		CreatedBy: uuid.MustParse(inp.CreatedBy),
	})
	if err != nil {
		return repository.CreateRoomOutput{}, fmt.Errorf("failed to query: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return repository.CreateRoomOutput{}, fmt.Errorf("failed to commit: %w", err)
	}

	return repository.CreateRoomOutput{}, nil
}

var _ repository.RoomRepository = new(RoomRepositoryOnDB)
