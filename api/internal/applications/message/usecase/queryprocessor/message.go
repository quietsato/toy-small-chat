package queryprocessor

type GetMessagesInput struct {
	RoomID string
}

type GetMessagesOutput struct {
	Messages []string
}

type Message struct {
	ID        string
	Author    string
	Content   string
	CreatedAt string
}

type MessageQueryProcessor interface {
	GetMessages(roomID string) ([]Message, error)
}
