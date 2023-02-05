package entity

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"time"
)

type Message struct {
	Body      string    `json:"body" vall:"Сообщение"`
	Sender    FindUser  `json:"sender"`
	Recipient FindUser  `json:"recipient"`
	Date      time.Time `json:"date"`
}

type FindMessage struct {
	Body string    `json:"body"`
	Date time.Time `json:"date"`
}

type ListMessage []FindMessage

func (a Message) Validate() error {
	return validation.ValidateStruct(&a,
		validation.Field(&a.Body,
			validation.Required),
	)
}
