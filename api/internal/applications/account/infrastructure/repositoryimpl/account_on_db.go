package repositoryimpl

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/quietsato/toy-small-chat/api/internal/applications/account/usecase/repository"
	"github.com/quietsato/toy-small-chat/api/internal/db"
)

func NewAccountRepositoryOnDB(pool *pgxpool.Pool) *AccountRepositoryOnDB {
	return &AccountRepositoryOnDB{pool}
}

type AccountRepositoryOnDB struct {
	pool *pgxpool.Pool
}

func (r *AccountRepositoryOnDB) CreateAccount(ctx context.Context, inp repository.CreateAccountInput) (repository.CreateAccountOutput, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return repository.CreateAccountOutput{}, fmt.Errorf("failed to begin tx: %w", err)
	}
	defer func() {
		if err := tx.Rollback(ctx); err != nil {
			slog.ErrorContext(ctx, "failed to rollback", slog.Any("err", err))
		}
	}()

	queries := db.New(r.pool).WithTx(tx)
	res, err := queries.CreateAccount(ctx, db.CreateAccountParams{
		Username:     inp.UserName,
		PasswordHash: inp.PasswordHash,
	})
	if err != nil {
		return repository.CreateAccountOutput{}, fmt.Errorf("failed to create account: %w", err)
	}
	if err := tx.Commit(ctx); err != nil {
		return repository.CreateAccountOutput{}, fmt.Errorf("failed to commit: %w", err)
	}

	return repository.CreateAccountOutput{
		AccountID: res.String(),
	}, nil
}

var _ repository.AccountRepository = new(AccountRepositoryOnDB)
