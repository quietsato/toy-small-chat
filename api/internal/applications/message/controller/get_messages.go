package controller

import (
	"github.com/quietsato/toy-small-chat/api/internal/applications/message/usecase/queryprocessor"
)

type GetMessagesController struct {
	query queryprocessor.MessageQueryProcessor
}

func NewGetMessagesController(query queryprocessor.MessageQueryProcessor) *GetMessagesController {
	return &GetMessagesController{query}
}

func (c *GetMessagesController) GetMessages(inp GetMessagesInput) (GetMessagesOutput, error) {
	queryResult, err := c.query.GetMessages(inp.RoomID)
	if err != nil {
		return GetMessagesOutput{}, err
	}

	msgs := make([]Message, 0, len(queryResult))
	for _, msg := range queryResult {
		msgs = append(msgs, Message{
			ID:        msg.ID,
			Content:   msg.Content,
			Author:    msg.Author,
			CreatedAt: msg.CreatedAt,
		})
	}

	return GetMessagesOutput{
		Messages: msgs,
	}, nil
}

type GetMessagesInput struct {
	RoomID string `json:"roomId"`
}
type GetMessagesOutput struct {
	Messages []Message `json:"messages"`
}

type Message struct {
	ID        string `json:"id"`
	Content   string `json:"content"`
	Author    string `json:"author"`
	CreatedAt string `json:"createdAt"`
}
