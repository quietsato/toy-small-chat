package domain

import (
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type MessageID struct {
	uuid uuid.UUID
}

func ParseMessageID(s string) (MessageID, error) {
	u, err := uuid.Parse(s)
	if err != nil {
		return MessageID{}, fmt.Errorf("failed to parse uuid: %w", err)
	}
	return MessageID{uuid: u}, nil
}

func MessageIDFromUuid(u uuid.UUID) MessageID {
	return MessageID{uuid: u}
}

func (m MessageID) String() string {
	return m.uuid.String()
}

type MessageContent struct {
	content string
}

const (
	messageContentMinLength = 1
	messageContentMaxLength = 1000
)

var (
	ErrInvalidMessageContent = errors.New("invalid message content")
)

func NewMessageContent(s string) (MessageContent, error) {
	if len(s) < messageContentMinLength || len(s) > messageContentMaxLength {
		return MessageContent{}, ErrInvalidMessageContent
	}
	return MessageContent{content: s}, nil
}

func (m MessageContent) String() string {
	return m.content
}

type Message struct {
	id        MessageID
	roomID    RoomID
	senderID  AccountID
	content   MessageContent
	createdAt time.Time
}

func NewMessage(id MessageID, roomID RoomID, senderID AccountID, content MessageContent, createdAt time.Time) Message {
	return Message{
		id:        id,
		roomID:    roomID,
		senderID:  senderID,
		content:   content,
		createdAt: createdAt,
	}
}

func (m Message) ID() MessageID {
	return m.id
}

func (m Message) RoomID() RoomID {
	return m.roomID
}

func (m Message) SenderID() AccountID {
	return m.senderID
}

func (m Message) Content() MessageContent {
	return m.content
}

func (m Message) CreatedAt() time.Time {
	return m.createdAt
}

type Messages struct {
	messages []Message
}

func NewMessages(messages []Message) Messages {
	return Messages{messages: messages}
}

func (m Messages) List() []Message {
	return m.messages
}
