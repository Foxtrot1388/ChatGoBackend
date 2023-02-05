package v1

import (
	"ChatGo/internal/domain/entity"
	auth "ChatGo/server/midleware"
	"context"
)

type UseCaseUser interface {
	Create(user *entity.User) error
	Login(user *entity.User) (*entity.FindUser, error)
	Find(user string) (*entity.ListUser, error)
	AddContact(curuser *entity.FindUser, adduser *entity.FindUser) error
}

type ControllerUser struct {
	userUseCase UseCaseUser
}

func NewUserUseCase(userUseCase UseCaseUser) *ControllerUser {
	return &ControllerUser{userUseCase: userUseCase}
}

func (c *ControllerUser) Create(ctx context.Context, user *entity.User) error {
	return c.userUseCase.Create(user)
}

func (c *ControllerUser) Login(ctx context.Context, user *entity.User) (token string, err error) {

	result, err := c.userUseCase.Login(user)
	if err != nil {
		return "", err
	}

	token, err = auth.CreateJWT(result.Login)
	if err != nil {
		return "", err
	}
	return token, nil

}

func (c *ControllerUser) Find(ctx context.Context, user string) (*entity.ListUser, error) {
	return c.userUseCase.Find(user)
}

func (c *ControllerUser) AddContact(ctx context.Context, adduser *entity.FindUser) error {
	return c.userUseCase.AddContact(
		&entity.FindUser{
			Login: ctx.Value("User").(string),
		},
		adduser,
	)
}
