package domain_test

import (
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/quietsato/toy-small-chat/api/internal/domain"
	"github.com/stretchr/testify/require"
)

func TestParseMessageID(t *testing.T) {
	t.Run("valid UUID", func(t *testing.T) {
		validUUID := "550e8400-e29b-41d4-a716-446655440000"
		messageID, err := domain.ParseMessageID(validUUID)
		require.NoError(t, err)
		require.Equal(t, validUUID, messageID.String())
	})

	t.Run("invalid UUID", func(t *testing.T) {
		_, err := domain.ParseMessageID("not-a-uuid")
		require.Error(t, err)
	})

	t.Run("empty string", func(t *testing.T) {
		_, err := domain.ParseMessageID("")
		require.Error(t, err)
	})
}

func TestMessageIDFromUuid(t *testing.T) {
	u := uuid.New()
	messageID := domain.MessageIDFromUuid(u)
	require.Equal(t, u.String(), messageID.String())
}

func TestNewMessageContent(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		wantError bool
	}{
		{"valid single char", "a", false},
		{"valid message", "Hello, world!", false},
		{"valid with newline", "Hello\nworld", false},
		{"valid with unicode", "こんにちは", false},
		{"valid with emoji", "Hello 😀", false},
		{"valid max length 1000", strings.Repeat("a", 1000), false},
		{"empty string", "", true},
		{"too long 1001 chars", strings.Repeat("a", 1001), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			content, err := domain.NewMessageContent(tt.input)
			if tt.wantError {
				require.ErrorIs(t, err, domain.ErrInvalidMessageContent)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.input, content.String())
			}
		})
	}
}

func TestNewMessage(t *testing.T) {
	messageID := domain.MessageIDFromUuid(uuid.New())
	roomID := domain.RoomIDFromUuid(uuid.New())
	senderID := domain.AccountIDFromUuid(uuid.New())
	content, _ := domain.NewMessageContent("Hello, world!")
	now := time.Now()

	message := domain.NewMessage(messageID, roomID, senderID, content, now)

	require.Equal(t, messageID.String(), message.ID().String())
	require.Equal(t, roomID.String(), message.RoomID().String())
	require.Equal(t, senderID.String(), message.SenderID().String())
	require.Equal(t, content.String(), message.Content().String())
	require.True(t, message.CreatedAt().Equal(now))
}

func TestNewMessages(t *testing.T) {
	t.Run("with messages", func(t *testing.T) {
		content1, _ := domain.NewMessageContent("Hello")
		content2, _ := domain.NewMessageContent("World")
		msg1 := domain.NewMessage(
			domain.MessageIDFromUuid(uuid.New()),
			domain.RoomIDFromUuid(uuid.New()),
			domain.AccountIDFromUuid(uuid.New()),
			content1,
			time.Now(),
		)
		msg2 := domain.NewMessage(
			domain.MessageIDFromUuid(uuid.New()),
			domain.RoomIDFromUuid(uuid.New()),
			domain.AccountIDFromUuid(uuid.New()),
			content2,
			time.Now(),
		)

		messages := domain.NewMessages([]domain.Message{msg1, msg2})
		require.Len(t, messages.List(), 2)
	})

	t.Run("empty", func(t *testing.T) {
		messages := domain.NewMessages([]domain.Message{})
		require.Empty(t, messages.List())
	})
}
