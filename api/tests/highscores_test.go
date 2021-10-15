package tests

import (
	"bytes"
	"crate-run-api/middlewares"
	"crate-run-api/models"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-playground/assert/v2"
)

func TestGetGlobalHighScores(t *testing.T) {
	a := createTestApp()

	users := [...]*models.User{{Username: "one"}, {Username: "two"}, {Username: "three"}}

	for _, user := range users {
		a.DB.Create(user)
	}

	now := time.Now()
	a.DB.Create(&models.HighScore{
		UserID:         users[2].ID,
		LevelNumber:    1,
		CompletionTime: 10,
		Moves:          10,
		AchievedAt:     now,
	})

	a.DB.Create(&models.HighScore{
		UserID:         users[1].ID,
		LevelNumber:    1,
		CompletionTime: 10,
		Moves:          5,
		AchievedAt:     now,
	})

	a.DB.Create(&models.HighScore{
		UserID:         users[2].ID,
		LevelNumber:    2,
		CompletionTime: 10,
		Moves:          10,
		AchievedAt:     now,
	})

	a.DB.Create(&models.HighScore{
		UserID:         users[0].ID,
		LevelNumber:    3,
		CompletionTime: 10,
		Moves:          10,
		AchievedAt:     now.Add(-5 * time.Minute),
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/highscores", nil)
	a.Engine.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	var scores []models.UserTotalScore
	assert.Equal(t, json.Unmarshal(w.Body.Bytes(), &scores), nil)

	assert.Equal(t, len(scores), 3)
	assert.Equal(t, scores[0].User.ID, users[2].ID)
	assert.Equal(t, scores[1].User.ID, users[0].ID)
	assert.Equal(t, scores[2].User.ID, users[1].ID)
	assert.Equal(t, scores[0].TotalScore, uint(19))
	assert.Equal(t, scores[1].TotalScore, uint(10))
	assert.Equal(t, scores[2].TotalScore, uint(10))
}

func TestGetLevelHighScoreInvalidNumber(t *testing.T) {
	a := createTestApp()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/highscores/level/a", nil)
	a.Engine.ServeHTTP(w, req)

	assert.Equal(t, 400, w.Code)
}

func TestGetLevelHighScoreEmpty(t *testing.T) {
	a := createTestApp()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/highscores/level/1", nil)
	a.Engine.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	var scores []models.HighScore
	assert.Equal(t, json.Unmarshal(w.Body.Bytes(), &scores), nil)

	assert.Equal(t, len(scores), 0)
}

func TestGetLevelHighScore(t *testing.T) {
	a := createTestApp()

	users := [...]*models.User{
		{Username: "a"},
		{Username: "b"},
		{Username: "c"},
		{Username: "d"},
		{Username: "e"},
	}

	for _, user := range users {
		a.DB.Create(user)
	}

	now := time.Now()
	a.DB.Create(&models.HighScore{
		UserID:         users[0].ID,
		LevelNumber:    1,
		CompletionTime: 10,
		Moves:          10,
		AchievedAt:     now,
	})

	a.DB.Create(&models.HighScore{
		UserID:         users[1].ID,
		LevelNumber:    1,
		CompletionTime: 10,
		Moves:          5,
		AchievedAt:     now,
	})

	a.DB.Create(&models.HighScore{
		UserID:         users[2].ID,
		LevelNumber:    1,
		CompletionTime: 8,
		Moves:          5,
		AchievedAt:     now,
	})

	a.DB.Create(&models.HighScore{
		UserID:         users[3].ID,
		LevelNumber:    1,
		CompletionTime: 8,
		Moves:          5,
		AchievedAt:     now.Add(-5 * time.Minute),
	})

	a.DB.Create(&models.HighScore{
		UserID:         users[4].ID,
		LevelNumber:    2,
		CompletionTime: 2,
		Moves:          2,
		AchievedAt:     now,
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/highscores/level/1", nil)
	a.Engine.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	var scores []models.HighScore
	assert.Equal(t, json.Unmarshal(w.Body.Bytes(), &scores), nil)

	assert.Equal(t, len(scores), 4)
	assert.Equal(t, scores[0].Moves, uint(5))
	assert.Equal(t, scores[1].Moves, uint(5))
	assert.Equal(t, scores[2].Moves, uint(5))
	assert.Equal(t, scores[3].Moves, uint(10))
	assert.Equal(t, scores[0].User.ID, users[3].ID)
	assert.Equal(t, scores[1].User.ID, users[2].ID)
	assert.Equal(t, scores[2].User.ID, users[1].ID)
	assert.Equal(t, scores[3].User.ID, users[0].ID)
}

func TestGetUserHighScoresUnauthorized(t *testing.T) {
	a := createTestApp()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/highscores/user", nil)
	a.Engine.ServeHTTP(w, req)

	assert.Equal(t, 401, w.Code)
}

func TestGetUserHighScores(t *testing.T) {
	a := createTestApp()

	users := [...]*models.User{{Username: "a"}}

	for _, user := range users {
		a.DB.Create(user)
	}

	now := time.Now()
	a.DB.Create(&models.HighScore{
		UserID:         users[0].ID,
		LevelNumber:    1,
		CompletionTime: 10,
		Moves:          10,
		AchievedAt:     now,
	})

	a.DB.Create(&models.HighScore{
		UserID:         users[0].ID,
		LevelNumber:    2,
		CompletionTime: 10,
		Moves:          5,
		AchievedAt:     now,
	})

	a.DB.Create(&models.HighScore{
		UserID:         users[0].ID,
		LevelNumber:    3,
		CompletionTime: 8,
		Moves:          5,
		AchievedAt:     now,
	})

	token := a.TokenStore.CreateToken("a")

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/highscores/user", nil)
	req.Header.Add(middlewares.AuthTokenHeader, token)
	a.Engine.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	var scores map[uint]models.HighScore
	assert.Equal(t, json.Unmarshal(w.Body.Bytes(), &scores), nil)

	assert.Equal(t, len(scores), 3)

	h1, ok1 := scores[1]
	assert.Equal(t, ok1, true)
	assert.Equal(t, h1.Moves, uint(10))
	assert.Equal(t, h1.CompletionTime, uint(10))

	h2, ok2 := scores[2]
	assert.Equal(t, ok2, true)
	assert.Equal(t, h2.Moves, uint(5))
	assert.Equal(t, h2.CompletionTime, uint(10))

	h3, ok3 := scores[3]
	assert.Equal(t, ok3, true)
	assert.Equal(t, h3.Moves, uint(5))
	assert.Equal(t, h3.CompletionTime, uint(8))
}

func TestRegisterHighScoreUnauthorized(t *testing.T) {
	a := createTestApp()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/highscores/level/1", nil)
	a.Engine.ServeHTTP(w, req)

	assert.Equal(t, 401, w.Code)
}

func TestRegisterHighScoreInvalidLevel(t *testing.T) {
	a := createTestApp()

	a.DB.Create(&models.User{Username: "a"})
	token := a.TokenStore.CreateToken("a")

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/highscores/level/a", nil)
	req.Header.Add(middlewares.AuthTokenHeader, token)
	a.Engine.ServeHTTP(w, req)

	assert.Equal(t, 400, w.Code)
}

func TestRegisterHighScoreLevelNumberLimits(t *testing.T) {
	a := createTestApp()

	a.DB.Create(&models.User{Username: "a"})
	token := a.TokenStore.CreateToken("a")
	data, _ := json.Marshal(models.HighScore{Moves: 1, CompletionTime: 1})

	{
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("PUT", "/highscores/level/0", bytes.NewReader(data))
		req.Header.Add(middlewares.AuthTokenHeader, token)
		a.Engine.ServeHTTP(w, req)

		assert.Equal(t, 400, w.Code)
	}

	{
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("PUT", "/highscores/level/51", bytes.NewReader(data))
		req.Header.Add(middlewares.AuthTokenHeader, token)
		a.Engine.ServeHTTP(w, req)

		assert.Equal(t, 400, w.Code)
	}

	{
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("PUT", "/highscores/level/1", bytes.NewReader(data))
		req.Header.Add(middlewares.AuthTokenHeader, token)
		a.Engine.ServeHTTP(w, req)

		assert.Equal(t, 200, w.Code)
	}

	{
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("PUT", "/highscores/level/50", bytes.NewReader(data))
		req.Header.Add(middlewares.AuthTokenHeader, token)
		a.Engine.ServeHTTP(w, req)

		assert.Equal(t, 200, w.Code)
	}
}

func TestRegisterHighScoreCreateScore(t *testing.T) {
	a := createTestApp()

	user := models.User{Username: "a"}
	a.DB.Create(&user)
	token := a.TokenStore.CreateToken("a")
	data, _ := json.Marshal(models.HighScore{Moves: 4, CompletionTime: 5})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/highscores/level/1", bytes.NewReader(data))
	req.Header.Add(middlewares.AuthTokenHeader, token)
	a.Engine.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	var score models.HighScore
	a.DB.Model(&models.HighScore{}).Where(&models.HighScore{UserID: user.ID}).First(&score)

	assert.Equal(t, score.Moves, uint(4))
	assert.Equal(t, score.CompletionTime, uint(5))
	assert.Equal(t, score.LevelNumber, uint(1))
	assert.Equal(t, score.AchievedAt.Unix() > 0, true)
}

func TestRegisterHighScoreUpdateScore(t *testing.T) {
	a := createTestApp()

	user := models.User{Username: "a"}
	a.DB.Create(&user)
	token := a.TokenStore.CreateToken("a")
	now := time.Now()
	regScore := models.HighScore{
		UserID:         user.ID,
		Moves:          6,
		CompletionTime: 7,
		AchievedAt:     now,
		LevelNumber:    2,
	}
	a.DB.Create(&regScore)

	data, _ := json.Marshal(models.HighScore{Moves: 4, CompletionTime: 5})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/highscores/level/2", bytes.NewReader(data))
	req.Header.Add(middlewares.AuthTokenHeader, token)
	a.Engine.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	var score models.HighScore
	a.DB.Model(&models.HighScore{}).Where(&models.HighScore{UserID: user.ID}).First(&score)

	assert.Equal(t, score.Moves, uint(4))
	assert.Equal(t, score.CompletionTime, uint(5))
	assert.Equal(t, score.LevelNumber, uint(2))
	assert.Equal(t, score.AchievedAt.Unix() >= now.Unix(), true)
}

func TestRegisterHighScoreNoOp(t *testing.T) {
	a := createTestApp()

	user := models.User{Username: "a"}
	a.DB.Create(&user)
	token := a.TokenStore.CreateToken("a")
	now := time.Now()
	regScore := models.HighScore{
		UserID:         user.ID,
		Moves:          6,
		CompletionTime: 7,
		AchievedAt:     now,
		LevelNumber:    2,
	}
	a.DB.Create(&regScore)

	data, _ := json.Marshal(models.HighScore{Moves: 8, CompletionTime: 9})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/highscores/level/2", bytes.NewReader(data))
	req.Header.Add(middlewares.AuthTokenHeader, token)
	a.Engine.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	var score models.HighScore
	a.DB.Model(&models.HighScore{}).Where(&models.HighScore{UserID: user.ID}).First(&score)

	assert.Equal(t, score.Moves, uint(6))
	assert.Equal(t, score.CompletionTime, uint(7))
	assert.Equal(t, score.LevelNumber, uint(2))
	assert.Equal(t, score.AchievedAt.Unix() == now.Unix(), true)
}
