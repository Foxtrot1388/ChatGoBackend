package message_usecase

import (
	"ChatGo/internal/domain/entity"
	"time"
)

type Repository interface {
	CreateMessage(mes *entity.Message) (string, error)
	FindOne(user string) (*entity.FindUser, error)
	ListMessages(sender *entity.FindUser, recipient *entity.FindUser, afterAt interface{}) (*entity.ListMessage, error)
}

type UseCase struct {
	repo Repository
}

func New(r Repository) *UseCase {
	return &UseCase{
		repo: r,
	}
}

func (r *UseCase) ListMessages(recipient string, sender string, afterAt interface{}) (*entity.ListMessage, error) {
	return r.repo.ListMessages(&entity.FindUser{Login: sender}, &entity.FindUser{Login: recipient}, afterAt)
}

func (r *UseCase) CreateMessage(body string, recipient string, sender string) (string, error) {

	adduser, err := r.repo.FindOne(recipient)
	if err != nil {
		return "", err
	}

	Message := entity.Message{
		Body:      body,
		Recipient: *adduser,
		Date:      time.Now(),
		Sender:    entity.FindUser{Login: sender},
	}

	err = Message.Validate()
	if err != nil {
		return "", err
	}

	return r.repo.CreateMessage(&Message)

}
