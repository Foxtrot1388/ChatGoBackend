package entity

import (
	"ChatGo/config"
	"crypto/md5"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

type User struct {
	Login string `json:"login" vall:"Логин"`
	Pass  string `json:"pass" vall:"Пароль"`
}

type FindUser struct {
	Login string `json:"login" bson:"_id"`
}

type ListUser []FindUser

func (a User) GetHash() []byte {
	M5 := md5.New()
	M5.Write([]byte(a.Pass))
	M5.Write([]byte(config.Get().Salt))
	return M5.Sum(nil)
}

func (a FindUser) Validate() error {
	return validation.ValidateStruct(&a,
		validation.Field(&a.Login,
			validation.Required),
	)
}

func (a User) Validate() error {
	return validation.ValidateStruct(&a,
		validation.Field(&a.Login,
			validation.Required,
			validation.Length(4, 16).Error("Длинна должна быть от 4 до 16 символов"),
			is.UTFLetterNumeric.Error("Разрешенны только символы и цифры")),
		validation.Field(&a.Pass,
			validation.Required,
			validation.Length(8, 20).Error("Длинна должна быть от 8 до 20 символов"),
			is.UTFLetterNumeric.Error("Разрешенны только символы и цифры")),
	)
}
