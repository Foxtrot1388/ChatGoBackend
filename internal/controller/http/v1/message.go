package v1

import (
	"ChatGo/internal/domain/entity"
	"context"
)

type UseCaseMessage interface {
	CreateMessage(body string, recipient string, sender string) error
	ListMessages(recipient string, sender string, afterAt interface{}) (*entity.ListMessage, error)
}

type ControllerMessage struct {
	messageUseCase UseCaseMessage
}

func NewMessageUseCase(messageUseCase UseCaseMessage) *ControllerMessage {
	return &ControllerMessage{messageUseCase: messageUseCase}
}

func (c *ControllerMessage) CreateMessage(ctx context.Context, body string, recipient string) error {
	return c.messageUseCase.CreateMessage(body, recipient, ctx.Value("User").(string))
}

func (c *ControllerMessage) ListMessages(ctx context.Context, recipient string, afterAt interface{}) (*entity.ListMessage, error) {
	return c.messageUseCase.ListMessages(recipient, ctx.Value("User").(string), afterAt)
}
