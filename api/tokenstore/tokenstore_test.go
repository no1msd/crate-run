package tokenstore

import (
	"errors"
	"testing"
	"time"

	"github.com/go-playground/assert/v2"
)

type mockTimeProvider struct {
	currentTime *time.Time
}

func (m mockTimeProvider) Now() time.Time {
	return *m.currentTime
}

func (m mockTimeProvider) Since(t time.Time) time.Duration {
	return m.Now().Sub(t)
}

func TestCreateToken(t *testing.T) {
	store := New()

	token := store.CreateToken("foo")

	assert.NotEqual(t, token, "")
	assert.Equal(t, store.IsTokenValid(token), true)
	assert.Equal(t, store.IsTokenValid("notatoken"), false)

	username, err := store.GetUsernameForToken(token)
	assert.Equal(t, err, nil)
	assert.Equal(t, username, "foo")
}

func TestInvalidateToken(t *testing.T) {
	store := New()

	token := store.CreateToken("foo")

	assert.Equal(t, store.IsTokenValid(token), true)

	store.InvalidateToken(token)

	assert.Equal(t, store.IsTokenValid(token), false)

	username, err := store.GetUsernameForToken(token)
	assert.Equal(t, errors.Is(err, ErrInvalidToken), true)
	assert.Equal(t, username, "")
}

func TestTokenTimedOut(t *testing.T) {
	wayback := time.Date(2020, time.January, 1, 0, 0, 0, 0, time.UTC)

	store := New()
	store.TimeProvider = mockTimeProvider{currentTime: &wayback}
	token := store.CreateToken("foo")
	assert.Equal(t, store.IsTokenValid(token), true)

	wayback = wayback.Add(store.TokenTimeout)

	assert.Equal(t, store.IsTokenValid(token), true)

	wayback = wayback.Add(time.Microsecond)

	assert.Equal(t, store.IsTokenValid(token), false)
}

func TestTokenResetTimeout(t *testing.T) {
	wayback := time.Date(2020, time.January, 1, 0, 0, 0, 0, time.UTC)

	store := New()
	store.TimeProvider = mockTimeProvider{currentTime: &wayback}
	token := store.CreateToken("foo")
	assert.Equal(t, store.IsTokenValid(token), true)

	wayback = wayback.Add(time.Minute)

	assert.Equal(t, store.IsTokenValid(token), true)

	err := store.ResetTokenTimeout(token)

	assert.Equal(t, err, nil)
	assert.Equal(t, store.IsTokenValid(token), true)

	wayback = wayback.Add(store.TokenTimeout)

	assert.Equal(t, store.IsTokenValid(token), true)

	wayback = wayback.Add(time.Microsecond)

	assert.Equal(t, store.IsTokenValid(token), false)
}

func TestTokenResetInvalidToken(t *testing.T) {
	store := New()

	err := store.ResetTokenTimeout("notatoken")

	assert.Equal(t, errors.Is(err, ErrInvalidToken), true)
	assert.Equal(t, store.IsTokenValid("notatoken"), false)
}
