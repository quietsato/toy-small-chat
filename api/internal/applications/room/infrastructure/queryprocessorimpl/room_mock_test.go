package queryprocessorimpl_test

import (
	"testing"

	"github.com/quietsato/toy-small-chat/api/internal/applications/room/infrastructure/queryprocessorimpl"
	"github.com/quietsato/toy-small-chat/api/internal/applications/room/usecase/queryprocessor"
	"github.com/stretchr/testify/require"
)

func TestMockRoomQueryProcessor_GetRooms(t *testing.T) {
	t.Parallel()

	t.Run("空のルーム一覧を返す", func(t *testing.T) {
		t.Parallel()

		mock := &queryprocessorimpl.MockRoomQueryProcessor{}

		out, err := mock.GetRooms(t.Context(), queryprocessor.GetRoomsInput{})

		require.NoError(t, err)
		require.Empty(t, out.Rooms)
	})
}
