package usecase

type CreateMessage struct{}

type CreateMessageInput struct {
}

type CreateMessageOutput struct{}

func (c *CreateMessageOutput) Error() string {
	panic("unimplemented")
}

var _ error = new(CreateMessageOutput)

type CreateMessageError struct{}

func (uc *CreateMessage) Execute(inp CreateMessageInput) (*CreateMessageOutput, error) {
	return nil, nil
}
