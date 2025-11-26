package queryprocessorimpl

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/quietsato/toy-small-chat/api/internal/applications/room/usecase/queryprocessor"
	"github.com/quietsato/toy-small-chat/api/internal/db"
)

type RoomQueryProcessorOnDB struct {
	queries *db.Queries
}

func NewRoomQueryProcessorOnDB(pool *pgxpool.Pool) *RoomQueryProcessorOnDB {
	return &RoomQueryProcessorOnDB{
		queries: db.New(pool),
	}
}

// GetRooms implements queryprocessor.RoomQueryProcessor.
func (r *RoomQueryProcessorOnDB) GetRooms(ctx context.Context, inp queryprocessor.GetRoomsInput) (queryprocessor.GetRoomsOutput, error) {
	rows, err := r.queries.GetRooms(ctx)
	if err != nil {
		return queryprocessor.GetRoomsOutput{}, fmt.Errorf("failed to get rooms: %w", err)
	}

	rooms := make([]queryprocessor.RoomDTO, len(rows))
	for i, row := range rows {
		rooms[i] = queryprocessor.RoomDTO{
			ID:        row.ID.String(),
			Name:      row.Name,
			CreatedBy: row.CreatedBy.String(),
			CreatedAt: row.CreatedAt.Time.Format(""),
			UpdatedAt: row.CreatedAt.Time.Format(""),
		}
	}

	return queryprocessor.GetRoomsOutput{
		Rooms: rooms,
	}, nil
}
