package main

import (
	storage "ChatGo/internal/adapters/db/mongodb"
	"ChatGo/internal/config"
	"ChatGo/internal/domain/entity"
	app "ChatGo/server"
	"context"
	"encoding/json"
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestAll(t *testing.T) {

	validation.ErrorTag = "vall"
	testlogin := "TestUser"
	pass := "12345678"
	var token string

	cfg := config.Get()
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(cfg.Mongo.URI))
	if err != nil {
		t.Error(err)
		return
	}
	repo := storage.New(client.Database("ChatGo"))

	t.Cleanup(func() {
		err := repo.Delete(&entity.User{Login: testlogin})
		if err != nil && err != mongo.ErrNoDocuments {
			t.Error(err)
			return
		}
	})

	t.Run("Password is incorrect", func(t *testing.T) {

		createanswer := app.Answer{
			Error: "Пароль: Длинна должна быть от 8 до 20 символов.",
			Data:  "",
		}

		handlerCreate := http.HandlerFunc(app.Create)
		rec := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/CreateUser", strings.NewReader(fmt.Sprintf("{\"Login\":\"%s\", \"Pass\":\"1\"}", testlogin)))
		handlerCreate.ServeHTTP(rec, req)

		var NewUserAnswer app.Answer
		err := json.NewDecoder(rec.Body).Decode(&NewUserAnswer)
		if err != nil {
			t.Error(err)
			return
		}

		assert.Equal(t, rec.Code, http.StatusBadRequest)
		assert.Equal(t, createanswer, NewUserAnswer)
	})

	t.Run("Login is incorrect", func(t *testing.T) {

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
	})

	t.Run("User is correct", func(t *testing.T) {

		createanswer := app.Answer{
			Error: "",
			Data:  "Ok",
		}

		handlerCreate := http.HandlerFunc(app.Create)
		rec := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/CreateUser", strings.NewReader(fmt.Sprintf("{\"Login\":\"%s\", \"Pass\":\"%s\"}", testlogin, pass)))
		handlerCreate.ServeHTTP(rec, req)

		var NewUserAnswer app.Answer
		err := json.NewDecoder(rec.Body).Decode(&NewUserAnswer)
		if err != nil {
			t.Error(err)
			return
		}

		assert.Equal(t, rec.Code, http.StatusOK)
		assert.Equal(t, createanswer, NewUserAnswer)
	})

	t.Run("Login failed", func(t *testing.T) {

		handlerCreate := http.HandlerFunc(app.Create)
		rec := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/LoginUser", strings.NewReader(fmt.Sprintf("{\"Login\":\"%s\", \"Pass\":\"87654321\"}", testlogin)))
		handlerCreate.ServeHTTP(rec, req)

		assert.Equal(t, rec.Code, http.StatusBadRequest)

	})

	t.Run("Login successful", func(t *testing.T) {

		handlerCreate := http.HandlerFunc(app.Create)
		rec := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/LoginUser", strings.NewReader(fmt.Sprintf("{\"Login\":\"%s\", \"Pass\":\"%s\"}", testlogin, pass)))
		handlerCreate.ServeHTTP(rec, req)

		assert.Equal(t, rec.Code, http.StatusOK)

		var NewLoginAnswer app.Answer
		err := json.NewDecoder(rec.Body).Decode(&NewLoginAnswer)
		if err != nil {
			t.Error(err)
			return
		}

		token = NewLoginAnswer.Data.(string)

	})

	t.Run("Add contact successful", func(t *testing.T) {

		handlerCreate := http.HandlerFunc(app.Create)
		rec := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/AddContact", strings.NewReader(fmt.Sprintf("{\"Login\":\"%s\"}", testlogin)))
		req.Header.Add("Authorization", "Bearer "+token)
		handlerCreate.ServeHTTP(rec, req)

		assert.Equal(t, rec.Code, http.StatusOK)

		addanswer := app.Answer{
			Error: "",
			Data:  "Ok",
		}

		var NewAddAnswer app.Answer
		err := json.NewDecoder(rec.Body).Decode(&NewAddAnswer)
		if err != nil {
			t.Error(err)
			return
		}

		assert.Equal(t, addanswer, NewAddAnswer)

	})

}
