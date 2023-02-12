package v1

import (
	"ChatGo/internal/domain/entity"
	auth "ChatGo/server/midleware"
	"context"
)

type UseCaseUser interface {
	CreateUser(user *entity.User) error
	LoginUser(user *entity.User) (*entity.FindUser, error)
	FindUser(user string) (*entity.ListUser, error)
}

type ControllerUser struct {
	userUseCase UseCaseUser
}

func NewUserUseCase(userUseCase UseCaseUser) *ControllerUser {
	return &ControllerUser{userUseCase: userUseCase}
}

func (c *ControllerUser) CreateUser(ctx context.Context, user *entity.User) error {
	return c.userUseCase.CreateUser(user)
}

func (c *ControllerUser) LoginUser(ctx context.Context, user *entity.User) (token string, err error) {

	result, err := c.userUseCase.LoginUser(user)
	if err != nil {
		return "", err
	}

	token, err = auth.CreateJWT(result.Login)
	if err != nil {
		return "", err
	}
	return token, nil

}

func (c *ControllerUser) FindUser(ctx context.Context, user string) (*entity.ListUser, error) {
	return c.userUseCase.FindUser(user)
}
