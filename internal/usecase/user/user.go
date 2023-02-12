package user_usecase

import (
	"ChatGo/internal/domain/entity"
)

type Repository interface {
	CreateUser(user *entity.User) error
	LoginUser(user *entity.User) (*entity.FindUser, error)
	FindUser(user string) (*entity.ListUser, error)
}

type UseCaseUser struct {
	repo Repository
}

func NewUserUseCase(r Repository) *UseCaseUser {
	return &UseCaseUser{
		repo: r,
	}
}

func (r *UseCaseUser) CreateUser(user *entity.User) error {

	err := user.Validate()
	if err != nil {
		return err
	}

	return r.repo.CreateUser(user)
}

func (r *UseCaseUser) LoginUser(user *entity.User) (*entity.FindUser, error) {
	return r.repo.LoginUser(user)
}

func (r *UseCaseUser) FindUser(user string) (*entity.ListUser, error) {
	return r.repo.FindUser(user)
}
