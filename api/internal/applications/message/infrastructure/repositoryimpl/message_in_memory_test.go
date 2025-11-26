package repositoryimpl_test

import (
	"testing"

	"github.com/quietsato/toy-small-chat/api/internal/applications/message/infrastructure/repositoryimpl"
	"github.com/quietsato/toy-small-chat/api/internal/applications/message/usecase/repository"
	"github.com/stretchr/testify/require"
)

func TestInMemoryMessageRepository_CreateMessage(t *testing.T) {
	t.Parallel()

	t.Run("メッセージが正常に作成される", func(t *testing.T) {
		t.Parallel()

		msgs := make([]repositoryimpl.Message, 0)
		repo := repositoryimpl.NewInMemoryMessageRepository(t.Context(), &msgs)

		err := repo.CreateMessage(t.Context(), repository.CreateMessageInput{
			AuthorID: "author-123",
			Content:  "Hello, World!",
			RoomID:   "room-456",
		})

		require.NoError(t, err)
		require.Len(t, msgs, 1)
		require.Equal(t, "Hello, World!", msgs[0].Content)
		require.NotEmpty(t, msgs[0].ID)
	})

	t.Run("複数メッセージが追加される", func(t *testing.T) {
		t.Parallel()

		msgs := make([]repositoryimpl.Message, 0)
		repo := repositoryimpl.NewInMemoryMessageRepository(t.Context(), &msgs)

		err := repo.CreateMessage(t.Context(), repository.CreateMessageInput{
			Content: "Message 1",
		})
		require.NoError(t, err)

		err = repo.CreateMessage(t.Context(), repository.CreateMessageInput{
			Content: "Message 2",
		})
		require.NoError(t, err)

		require.Len(t, msgs, 2)
	})
}

func TestNewInMemoryMessageRepository(t *testing.T) {
	t.Run("正しく初期化される", func(t *testing.T) {
		msgs := make([]repositoryimpl.Message, 0)
		repo := repositoryimpl.NewInMemoryMessageRepository(t.Context(), &msgs)

		require.NotNil(t, repo)
	})
}
