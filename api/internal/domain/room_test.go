package domain_test

import (
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/quietsato/toy-small-chat/api/internal/domain"
	"github.com/stretchr/testify/require"
)

// RoomID tests
func TestParseRoomID(t *testing.T) {
	t.Parallel()

	t.Run("valid UUID", func(t *testing.T) {
		t.Parallel()
		validUUID := "550e8400-e29b-41d4-a716-446655440000"
		roomID, err := domain.ParseRoomID(validUUID)
		require.NoError(t, err)
		require.Equal(t, validUUID, roomID.String())
	})

	t.Run("invalid UUID", func(t *testing.T) {
		t.Parallel()

		invalidUUID := "not-a-uuid"
		_, err := domain.ParseRoomID(invalidUUID)
		require.Error(t, err)
	})

	t.Run("empty string", func(t *testing.T) {
		t.Parallel()

		_, err := domain.ParseRoomID("")
		require.Error(t, err)
	})
}

func TestRoomIDFromUuid(t *testing.T) {
	t.Parallel()

	t.Run("create from UUID", func(t *testing.T) {
		t.Parallel()
		u := uuid.New()
		roomID := domain.RoomIDFromUuid(u)
		require.Equal(t, u.String(), roomID.String())
	})
}

func TestRoomID_String(t *testing.T) {
	t.Parallel()

	t.Run("convert to string", func(t *testing.T) {
		t.Parallel()
		validUUID := "550e8400-e29b-41d4-a716-446655440000"
		roomID, _ := domain.ParseRoomID(validUUID)
		require.Equal(t, validUUID, roomID.String())
	})
}

// RoomName tests
func TestNewRoomName(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		input     string
		expected  string
		wantError bool
	}{
		{"valid lowercase", "general", "general", false},
		{"valid uppercase", "GENERAL", "GENERAL", false},
		{"valid mixed case", "General", "General", false},
		{"valid with numbers", "room123", "room123", false},
		{"valid with space", "my room", "my room", false},
		{"valid with underscore", "my_room", "my_room", false},
		{"valid with hyphen", "my-room", "my-room", false},
		{"valid single char", "a", "a", false},
		{"valid max length 127", strings.Repeat("a", 127), strings.Repeat("a", 127), false},
		{"valid with special chars", "room@name", "room@name", false},
		{"valid with unicode", "room名前", "room名前", false},
		{"valid with emoji", "room 🎉", "room 🎉", false},
		// Trim and normalization tests
		{"trim leading space", "  room", "room", false},
		{"trim trailing space", "room  ", "room", false},
		{"trim both spaces", "  room  ", "room", false},
		{"newline converted to space", "room\nname", "room name", false},
		{"multiple newlines converted to single space", "room\n\nname", "room name", false},
		{"carriage return converted to space", "room\rname", "room name", false},
		{"crlf converted to space", "room\r\nname", "room name", false},
		{"trim after newline normalization", "\n room \n", "room", false},
		// Error cases
		{"empty string", "", "", true},
		{"only spaces", "   ", "", true},
		{"only newlines", "\n\n", "", true},
		{"too long 128 chars", strings.Repeat("a", 128), "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			roomName, err := domain.NewRoomName(tt.input)
			if tt.wantError {
				require.Error(t, err)
				require.ErrorIs(t, err, domain.ErrInvalidRoomName)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.expected, roomName.String())
			}
		})
	}
}

func TestRoomName_String(t *testing.T) {
	t.Parallel()

	t.Run("convert to string", func(t *testing.T) {
		t.Parallel()
		input := "test-room"
		roomName, _ := domain.NewRoomName(input)
		require.Equal(t, input, roomName.String())
	})
}

// Room tests
func TestNewRoom(t *testing.T) {
	t.Parallel()

	t.Run("create room", func(t *testing.T) {
		t.Parallel()
		roomID := domain.RoomIDFromUuid(uuid.New())
		roomName, _ := domain.NewRoomName("test-room")
		createdBy := domain.AccountIDFromUuid(uuid.New())
		now := time.Now()

		room := domain.NewRoom(roomID, roomName, createdBy, now, now)

		gotID := room.ID()
		require.Equal(t, roomID.String(), gotID.String())
		gotName := room.Name()
		require.Equal(t, roomName.String(), gotName.String())
		gotCreatedBy := room.CreatedBy()
		require.Equal(t, createdBy.String(), gotCreatedBy.String())
		require.True(t, room.CreatedAt().Equal(now))
		require.True(t, room.UpdatedAt().Equal(now))
	})
}

// Rooms tests
func TestNewRooms(t *testing.T) {
	t.Parallel()

	t.Run("create rooms list", func(t *testing.T) {
		t.Parallel()
		roomID1 := domain.RoomIDFromUuid(uuid.New())
		roomName1, _ := domain.NewRoomName("room1")
		roomID2 := domain.RoomIDFromUuid(uuid.New())
		roomName2, _ := domain.NewRoomName("room2")
		createdBy := domain.AccountIDFromUuid(uuid.New())
		now := time.Now()

		room1 := domain.NewRoom(roomID1, roomName1, createdBy, now, now)
		room2 := domain.NewRoom(roomID2, roomName2, createdBy, now, now)

		rooms := domain.NewRooms([]domain.Room{room1, room2})
		list := rooms.List()

		require.Len(t, list, 2)
	})

	t.Run("empty rooms list", func(t *testing.T) {
		t.Parallel()

		rooms := domain.NewRooms([]domain.Room{})
		list := rooms.List()

		require.Empty(t, list)
	})
}
