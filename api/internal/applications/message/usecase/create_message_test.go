package usecase_test

import (
	"testing"

	"github.com/quietsato/toy-small-chat/api/internal/applications/message/usecase"
	"github.com/stretchr/testify/require"
)

func TestCreateMessage_Execute(t *testing.T) {
	t.Parallel()

	t.Run("メッセージ作成が正常に完了する", func(t *testing.T) {
		t.Parallel()

		uc := &usecase.CreateMessage{}
		out, err := uc.Execute(usecase.CreateMessageInput{})

		require.NoError(t, err)
		require.Nil(t, out)
	})
}
