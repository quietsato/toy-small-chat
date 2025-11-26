package queryprocessorimpl

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/quietsato/toy-small-chat/api/internal/applications/account/usecase/queryprocessor"
	"github.com/quietsato/toy-small-chat/api/internal/db"
)

type AccountQueryProcessorOnDB struct {
	pool *pgxpool.Pool
}

func NewAccountQueryProcessorOnDB(pool *pgxpool.Pool) *AccountQueryProcessorOnDB {
	return &AccountQueryProcessorOnDB{pool}
}

func (q *AccountQueryProcessorOnDB) GetLoginCredential(ctx context.Context, inp queryprocessor.GetLoginCredentialInput) (queryprocessor.GetLoginCredentialOutput, error) {
	queries := db.New(q.pool)
	res, err := queries.GetLoginCredential(ctx, inp.UserName)
	// アカウントが見つからなかった場合
	if errors.Is(err, sql.ErrNoRows) {
		return queryprocessor.GetLoginCredentialOutput{}, queryprocessor.ErrAccountNotFound
	}
	// その他の DB 起因のエラー
	if err != nil {
		return queryprocessor.GetLoginCredentialOutput{}, fmt.Errorf("failed to query: %w", err)
	}

	return queryprocessor.GetLoginCredentialOutput{
		AccountID:    res.ID,
		PasswordHash: res.PasswordHash,
	}, nil

}

var _ queryprocessor.AccountQueryProcessor = new(AccountQueryProcessorOnDB)
