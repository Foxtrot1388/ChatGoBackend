package v1

import (
	"ChatGo/internal/domain/entity"
	"context"
)

type UseCaseContact interface {
	AddContact(curuser *entity.FindUser, adduser *entity.FindUser) (string, error)
	DeleteContact(id string) error
	ListContact(login string) (*entity.ListContact, error)
}

type ControllerContact struct {
	contactUseCase UseCaseContact
}

func NewContactUseCase(contactUseCase UseCaseContact) *ControllerContact {
	return &ControllerContact{contactUseCase: contactUseCase}
}

func (c *ControllerContact) AddContact(ctx context.Context, adduser *entity.FindUser) (string, error) {
	return c.contactUseCase.AddContact(
		&entity.FindUser{
			Login: ctx.Value("User").(string),
		},
		adduser,
	)
}

func (c *ControllerContact) DeleteContact(ctx context.Context, id string) error {
	return c.contactUseCase.DeleteContact(id)
}

func (c *ControllerContact) ListContact(ctx context.Context) (*entity.ListContact, error) {
	return c.contactUseCase.ListContact(ctx.Value("User").(string))
}
