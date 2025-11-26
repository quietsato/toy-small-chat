package domain

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
)

type RoomID struct {
	uuid uuid.UUID
}

func ParseRoomID(s string) (RoomID, error) {
	u, err := uuid.Parse(s)
	if err != nil {
		return RoomID{}, fmt.Errorf("failed to parse uuid: %w", err)
	}
	return RoomID{uuid: u}, nil
}

func RoomIDFromUuid(u uuid.UUID) RoomID {
	return RoomID{uuid: u}
}

func (r RoomID) String() string {
	return r.uuid.String()
}

type RoomName struct {
	name string
}

const (
	roomNameMinLength = 1
	roomNameMaxLength = 127
)

var (
	ErrInvalidRoomName = errors.New("invalid room name")
	newlineRegExp      = regexp.MustCompile(`[\r\n]+`)
)

func NewRoomName(s string) (RoomName, error) {
	// Replace newlines with spaces
	normalized := newlineRegExp.ReplaceAllString(s, " ")
	// Trim leading and trailing whitespace
	normalized = strings.TrimSpace(normalized)

	if len(normalized) < roomNameMinLength || len(normalized) > roomNameMaxLength {
		return RoomName{}, ErrInvalidRoomName
	}
	return RoomName{name: normalized}, nil
}

func (r RoomName) String() string {
	return r.name
}

type Room struct {
	id        RoomID
	name      RoomName
	createdBy AccountID
	createdAt time.Time
	updatedAt time.Time
}

func NewRoom(id RoomID, name RoomName, createdBy AccountID, createdAt, updatedAt time.Time) Room {
	return Room{
		id:        id,
		name:      name,
		createdBy: createdBy,
		createdAt: createdAt,
		updatedAt: updatedAt,
	}
}

func (r Room) ID() RoomID {
	return r.id
}

func (r Room) Name() RoomName {
	return r.name
}

func (r Room) CreatedBy() AccountID {
	return r.createdBy
}

func (r Room) CreatedAt() time.Time {
	return r.createdAt
}

func (r Room) UpdatedAt() time.Time {
	return r.updatedAt
}

type Rooms struct {
	rooms []Room
}

func NewRooms(rooms []Room) Rooms {
	return Rooms{rooms: rooms}
}

func (r Rooms) List() []Room {
	return r.rooms
}
