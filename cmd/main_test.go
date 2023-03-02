package main

import (
	storage "ChatGo/internal/adapters/db/mongodb"
	"ChatGo/internal/domain/entity"
	app "ChatGo/server"
	"context"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/steinfletcher/apitest"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
	"testing"
)

func init() {
	validation.ErrorTag = "vall"
}

func TestUserPasswordIsIncorrect(t *testing.T) {

	apitest.
		HandlerFunc(app.Create).
		Post("/CreateUser").
		Bodyf("{\"Login\":\"%s\", \"Pass\":\"1\"}", "TestLogin").
		Expect(t).
		Status(http.StatusBadRequest).
		Body(`{"error": "Пароль: Длинна должна быть от 8 до 20 символов.", "data": ""}`).
		End()

}

func TestUserLoginIsIncorrect(t *testing.T) {

	apitest.
		HandlerFunc(app.Create).
		Post("/CreateUser").
		Body("{\"Login\":\"Den*\", \"Pass\":\"12345678\"}").
		Expect(t).
		Status(http.StatusBadRequest).
		Body(`{"error": "Логин: Разрешенны только символы и цифры.", "data": ""}`).
		End()

}

func TestUserIsCorrect(t *testing.T) {

	testlogin := "TestUser1"

	t.Cleanup(func() {

		repo, err := storage.New(context.TODO())
		assert.Nil(t, err)

		err = repo.DeleteUser(&entity.User{Login: testlogin})
		if err != nil && err != mongo.ErrNoDocuments {
			assert.Nil(t, err)
		}

	})

	apitest.
		HandlerFunc(app.Create).
		Post("/CreateUser").
		Bodyf("{\"Login\":\"%s\", \"Pass\":\"%s\"}", testlogin, "12345678").
		Expect(t).
		Status(http.StatusOK).
		Body(`{"error": "", "data": "Ok"}`).
		End()

}

func TestLoginFailed(t *testing.T) {

	repo, err := storage.New(context.TODO())
	assert.Nil(t, err)

	testlogin := "TestUser2"
	pass := "12345678"

	err = repo.CreateUser(&entity.User{
		Login: testlogin,
		Pass:  pass,
	})
	assert.Nil(t, err)

	t.Cleanup(func() {

		err = repo.DeleteUser(&entity.User{Login: testlogin})
		if err != nil && err != mongo.ErrNoDocuments {
			assert.Nil(t, err)
		}

	})

	apitest.
		HandlerFunc(app.Login).
		Post("/LoginUser").
		Bodyf("{\"Login\":\"%s\", \"Pass\":\"87654321\"}", testlogin).
		Expect(t).
		Status(http.StatusBadRequest).
		End()

}

func TestLoginSuccessful(t *testing.T) {

	repo, err := storage.New(context.TODO())
	assert.Nil(t, err)

	testlogin := "TestUser3"
	pass := "12345678"

	err = repo.CreateUser(&entity.User{
		Login: testlogin,
		Pass:  pass,
	})
	assert.Nil(t, err)

	t.Cleanup(func() {

		err = repo.DeleteUser(&entity.User{Login: testlogin})
		if err != nil && err != mongo.ErrNoDocuments {
			assert.Nil(t, err)
		}

	})

	apitest.
		HandlerFunc(app.Login).
		Post("/LoginUser").
		Bodyf("{\"Login\":\"%s\", \"Pass\":\"%s\"}", testlogin, pass).
		Expect(t).
		Status(http.StatusOK).
		End()

}

func TestAddContactSuccessful(t *testing.T) {

	repo, err := storage.New(context.TODO())
	if err != nil {
		assert.Nil(t, err)
	}

	testlogin := "TestUser4"
	pass := "12345678"

	err = repo.CreateUser(&entity.User{
		Login: testlogin,
		Pass:  pass,
	})
	if err != nil {
		assert.Nil(t, err)
	}

	t.Cleanup(func() {

		err = repo.DeleteUser(&entity.User{Login: testlogin})
		if err != nil && err != mongo.ErrNoDocuments {
			assert.Nil(t, err)
		}

	})

	apitest.
		HandlerFunc(app.AddContact).
		Post("/AddContact").
		WithContext(context.WithValue(context.TODO(), "User", testlogin)).
		Bodyf("{\"Login\":\"%s\"}", testlogin).
		Expect(t).
		Status(http.StatusOK).
		End()

}
