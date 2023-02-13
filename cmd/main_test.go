package main

import (
	storage "ChatGo/internal/adapters/db/mongodb"
	"ChatGo/internal/domain/entity"
	app "ChatGo/server"
	"context"
	"encoding/json"
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func init() {
	validation.ErrorTag = "vall"
}

func TestUserPasswordIsIncorrect(t *testing.T) {

	createanswer := app.Answer{
		Error: "Пароль: Длинна должна быть от 8 до 20 символов.",
		Data:  "",
	}

	handlerCreate := http.HandlerFunc(app.Create)
	rec := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/CreateUser", strings.NewReader(fmt.Sprintf("{\"Login\":\"%s\", \"Pass\":\"1\"}", "TestLogin")))
	handlerCreate.ServeHTTP(rec, req)

	var NewUserAnswer app.Answer
	err := json.NewDecoder(rec.Body).Decode(&NewUserAnswer)
	if err != nil {
		t.Error(err)
		return
	}

	assert.Equal(t, rec.Code, http.StatusBadRequest)
	assert.Equal(t, createanswer, NewUserAnswer)

}

func TestUserLoginIsIncorrect(t *testing.T) {

	createanswer := app.Answer{
		Error: "Логин: Разрешенны только символы и цифры.",
		Data:  "",
	}

	handlerCreate := http.HandlerFunc(app.Create)
	rec := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/CreateUser", strings.NewReader("{\"Login\":\"Den*\", \"Pass\":\"12345678\"}"))
	handlerCreate.ServeHTTP(rec, req)

	var NewUserAnswer app.Answer
	err := json.NewDecoder(rec.Body).Decode(&NewUserAnswer)
	if err != nil {
		t.Error(err)
		return
	}

	assert.Equal(t, rec.Code, http.StatusBadRequest)
	assert.Equal(t, createanswer, NewUserAnswer)

}

func TestUserIsCorrect(t *testing.T) {

	testlogin := "TestUser1"

	t.Cleanup(func() {

		repo, err := storage.New(context.TODO())
		if err != nil {
			t.Error(err)
			return
		}

		err = repo.DeleteUser(&entity.User{Login: testlogin})
		if err != nil && err != mongo.ErrNoDocuments {
			t.Error(err)
			return
		}

	})

	createanswer := app.Answer{
		Error: "",
		Data:  "Ok",
	}

	handlerCreate := http.HandlerFunc(app.Create)
	rec := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/CreateUser", strings.NewReader(fmt.Sprintf("{\"Login\":\"%s\", \"Pass\":\"%s\"}", testlogin, "12345678")))
	handlerCreate.ServeHTTP(rec, req)

	var NewUserAnswer app.Answer
	err := json.NewDecoder(rec.Body).Decode(&NewUserAnswer)
	if err != nil {
		t.Error(err)
		return
	}

	assert.Equal(t, rec.Code, http.StatusOK)
	assert.Equal(t, createanswer, NewUserAnswer)

}

func TestLoginFailed(t *testing.T) {

	repo, err := storage.New(context.TODO())
	if err != nil {
		t.Error(err)
		return
	}

	testlogin := "TestUser2"
	pass := "12345678"

	err = repo.CreateUser(&entity.User{
		Login: testlogin,
		Pass:  pass,
	})
	if err != nil {
		t.Error(err)
		return
	}

	t.Cleanup(func() {

		err = repo.DeleteUser(&entity.User{Login: testlogin})
		if err != nil && err != mongo.ErrNoDocuments {
			t.Error(err)
			return
		}

	})

	handlerCreate := http.HandlerFunc(app.Login)
	rec := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/LoginUser", strings.NewReader(fmt.Sprintf("{\"Login\":\"%s\", \"Pass\":\"87654321\"}", testlogin)))
	handlerCreate.ServeHTTP(rec, req)

	assert.Equal(t, rec.Code, http.StatusBadRequest)

}

func TestLoginSuccessful(t *testing.T) {

	repo, err := storage.New(context.TODO())
	if err != nil {
		t.Error(err)
		return
	}

	testlogin := "TestUser3"
	pass := "12345678"

	err = repo.CreateUser(&entity.User{
		Login: testlogin,
		Pass:  pass,
	})
	if err != nil {
		t.Error(err)
		return
	}

	t.Cleanup(func() {

		err = repo.DeleteUser(&entity.User{Login: testlogin})
		if err != nil && err != mongo.ErrNoDocuments {
			t.Error(err)
			return
		}

	})

	handlerCreate := http.HandlerFunc(app.Login)
	rec := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/LoginUser", strings.NewReader(fmt.Sprintf("{\"Login\":\"%s\", \"Pass\":\"%s\"}", testlogin, pass)))
	handlerCreate.ServeHTTP(rec, req)

	assert.Equal(t, rec.Code, http.StatusOK)

	var NewLoginAnswer app.Answer
	err = json.NewDecoder(rec.Body).Decode(&NewLoginAnswer)
	if err != nil {
		t.Error(err)
		return
	}

}

func TestAddContactSuccessful(t *testing.T) {

	repo, err := storage.New(context.TODO())
	if err != nil {
		t.Error(err)
		return
	}

	testlogin := "TestUser4"
	pass := "12345678"

	err = repo.CreateUser(&entity.User{
		Login: testlogin,
		Pass:  pass,
	})
	if err != nil {
		t.Error(err)
		return
	}

	t.Cleanup(func() {

		err = repo.DeleteUser(&entity.User{Login: testlogin})
		if err != nil && err != mongo.ErrNoDocuments {
			t.Error(err)
			return
		}

	})

	handlerCreate := http.HandlerFunc(app.AddContact)
	rec := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/AddContact", strings.NewReader(fmt.Sprintf("{\"Login\":\"%s\"}", testlogin)))
	handlerCreate.ServeHTTP(rec, req.WithContext(context.WithValue(req.Context(), "User", testlogin)))

	assert.Equal(t, rec.Code, http.StatusOK)

}
