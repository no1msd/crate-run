package tests

import (
	"crate-run-api/app"
	"crate-run-api/middlewares"
	"crate-run-api/models"
	"crate-run-api/utils"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-playground/assert/v2"
)

func TestLoginSuccess(t *testing.T) {
	a := createTestApp()

	a.DB.Create(&models.User{
		Username:     "test",
		PasswordHash: utils.HashPassword("test"),
		Nickname:     "Test",
	})

	data := `{"username":"test", "password": "test"}`

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/session", strings.NewReader(data))
	a.Engine.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	var session app.Session
	assert.Equal(t, json.Unmarshal(w.Body.Bytes(), &session), nil)

	assert.Equal(t, session.User.ID > 0, true)
	assert.Equal(t, session.User.Username, "test")
	assert.Equal(t, session.User.Nickname, "Test")

	assert.Equal(t, session.Token != "", true)
	assert.Equal(t, a.TokenStore.IsTokenValid(session.Token), true)
	username, _ := a.TokenStore.GetUsernameForToken(session.Token)
	assert.Equal(t, username, "test")
}

func TestLoginFailed(t *testing.T) {
	a := createTestApp()

	a.DB.Create(&models.User{
		Username:     "test",
		PasswordHash: utils.HashPassword("test"),
		Nickname:     "Test",
	})

	data := `{"username":"test", "password": "tst"}`

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/session", strings.NewReader(data))
	a.Engine.ServeHTTP(w, req)

	assert.Equal(t, 401, w.Code)
}

func TestLogoutUnathorized(t *testing.T) {
	a := createTestApp()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/session", nil)
	a.Engine.ServeHTTP(w, req)

	assert.Equal(t, 401, w.Code)
}

func TestLogout(t *testing.T) {
	a := createTestApp()

	a.DB.Create(&models.User{
		Username: "test",
		Nickname: "Test",
	})

	token := a.TokenStore.CreateToken("test")

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/session", nil)
	req.Header.Add(middlewares.AuthTokenHeader, token)
	a.Engine.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	assert.Equal(t, a.TokenStore.IsTokenValid(token), false)
}
