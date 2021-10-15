package tests

import (
	"crate-run-api/app"
	"crate-run-api/middlewares"
	"crate-run-api/models"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-playground/assert/v2"
)

func TestUserDetailsUnathorized(t *testing.T) {
	a := createTestApp()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/user", nil)
	a.Engine.ServeHTTP(w, req)

	assert.Equal(t, 401, w.Code)
}

func TestUserDetails(t *testing.T) {
	a := createTestApp()

	a.DB.Create(&models.User{
		Username: "test",
		Nickname: "Test"},
	)

	token := a.TokenStore.CreateToken("test")

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/user", nil)
	req.Header.Add(middlewares.AuthTokenHeader, token)
	a.Engine.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	var user models.User
	assert.Equal(t, json.Unmarshal(w.Body.Bytes(), &user), nil)

	assert.Equal(t, user.ID > 0, true)
	assert.Equal(t, user.Username, "test")
	assert.Equal(t, user.Nickname, "Test")
}

func TestChangeNicknameUnathorized(t *testing.T) {
	a := createTestApp()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PATCH", "/user/nickname", strings.NewReader("Other"))
	a.Engine.ServeHTTP(w, req)

	assert.Equal(t, 401, w.Code)
}

func TestChangeNickname(t *testing.T) {
	a := createTestApp()

	a.DB.Create(&models.User{
		Username: "test",
		Nickname: "Test"},
	)

	token := a.TokenStore.CreateToken("test")

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PATCH", "/user/nickname", strings.NewReader("Other"))
	req.Header.Add(middlewares.AuthTokenHeader, token)
	a.Engine.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	var user models.User
	a.DB.Model(&models.User{}).Where(models.User{Username: "test"}).First(&user)
	assert.Equal(t, user.Nickname, "Other")
}

func TestRegisterUserConflict(t *testing.T) {
	a := createTestApp()

	a.DB.Create(&models.User{
		Username: "test",
		Nickname: "Test"},
	)

	data := `{"username":"test", "password": "test"}`

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/user", strings.NewReader(data))
	a.Engine.ServeHTTP(w, req)

	assert.Equal(t, 409, w.Code)
}

func TestRegisterUser(t *testing.T) {
	a := createTestApp()

	data := `{"password": "test"}`

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/user", strings.NewReader(data))
	a.Engine.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	var session app.Session
	assert.Equal(t, json.Unmarshal(w.Body.Bytes(), &session), nil)

	assert.Equal(t, session.User.ID > 0, true)
	assert.Equal(t, session.User.Username != "", true)
	assert.Equal(t, session.User.Nickname, "")

	assert.Equal(t, session.Token != "", true)
	assert.Equal(t, a.TokenStore.IsTokenValid(session.Token), true)
	username, _ := a.TokenStore.GetUsernameForToken(session.Token)
	assert.Equal(t, username, session.User.Username)
}
