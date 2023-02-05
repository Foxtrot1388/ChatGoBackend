package user_usecase

import (
	"ChatGo/internal/domain/entity"
)

type Repository interface {
	Create(user *entity.User) error
	Login(user *entity.User) (*entity.FindUser, error)
	Find(user string) (*entity.ListUser, error)
	FindOne(user string) (*entity.FindUser, error)
	AddContact(curuser *entity.FindUser, adduser *entity.FindUser) error
}

type UseCase struct {
	repo Repository
}

func New(r Repository) *UseCase {
	return &UseCase{
		repo: r,
	}
}

func (r *UseCase) Create(user *entity.User) error {

	err := user.Validate()
	if err != nil {
		return err
	}

	return r.repo.Create(user)
}

func (r *UseCase) Login(user *entity.User) (*entity.FindUser, error) {
	return r.repo.Login(user)
}

func (r *UseCase) Find(user string) (*entity.ListUser, error) {
	return r.repo.Find(user)
}

func (r *UseCase) AddContact(curuser *entity.FindUser, user *entity.FindUser) error {

	err := user.Validate()
	if err != nil {
		return err
	}

	adduser, err := r.repo.FindOne(user.Login)
	if err != nil {
		return err
	}

	return r.repo.AddContact(curuser, adduser)

}
