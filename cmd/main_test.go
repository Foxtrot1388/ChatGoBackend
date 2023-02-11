package main

import (
	app "ChatGo/server"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestAll(t *testing.T) {
	handlerCreate := http.HandlerFunc(app.Create)
	rec := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/CreateUser", strings.NewReader("{\n    \"Login\":\"Denis\",\n    \"Pass\":\"12345678\"\n}"))
	handlerCreate.ServeHTTP(rec, req)
	assert.Equal(t, rec.Code, http.StatusBadRequest)
	assert.Equal(t, []byte(""), rec.Body.Bytes())
}
