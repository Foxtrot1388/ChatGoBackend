package contact_usecase

import "ChatGo/internal/domain/entity"

type Repository interface {
	AddContact(curuser *entity.FindUser, adduser *entity.FindUser) (string, error)
	DeleteContact(id string) error
	ListContact(login string) (*entity.ListContact, error)
	FindOneUser(user string) (*entity.FindUser, error)
}

type UseCase struct {
	repo Repository
}

func New(r Repository) *UseCase {
	return &UseCase{
		repo: r,
	}
}

func (r *UseCase) ListContact(login string) (*entity.ListContact, error) {
	return r.repo.ListContact(login)
}

func (r *UseCase) AddContact(curuser *entity.FindUser, user *entity.FindUser) (string, error) {

	err := user.Validate()
	if err != nil {
		return "", err
	}

	adduser, err := r.repo.FindOneUser(user.Login)
	if err != nil {
		return "", err
	}

	return r.repo.AddContact(curuser, adduser)

}

func (r *UseCase) DeleteContact(id string) error {
	return r.repo.DeleteContact(id)
}
