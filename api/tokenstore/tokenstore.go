// Package tokenstore manages access tokens and their associated data
package tokenstore

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
)

// TimeProvider is an interface over the time package so it can be mocked from tests.
type TimeProvider interface {
	Now() time.Time
	Since(t time.Time) time.Duration
}

type realTimeProvider struct{}

func (realTimeProvider) Now() time.Time {
	return time.Now()
}

func (realTimeProvider) Since(t time.Time) time.Duration {
	return time.Since(t)
}

// ErrInvalidToken is an error returned when trying to access the data associated with an
// invalid token. Both unknown and timed out tokens are considered invalid.
var ErrInvalidToken = errors.New("invalid token")

type sessionData struct {
	username string
	updated  time.Time
}

// MemoryStore is a thread safe in-memory token store.
type MemoryStore struct {
	store map[string]*sessionData
	mut   sync.Mutex

	// TimeProvider holds the currently used implementation of the TimeProvider interface.
	// Tests can overwrite this variable to supply their own mock implementation.
	TimeProvider TimeProvider

	// TokenTimeout is the duration of the validity of an issued token.
	TokenTimeout time.Duration
}

// CreateToken returns a newly created token for the given username that is valid from the time
// of its creation until TokenTimeout amount of time passed.
func (s *MemoryStore) CreateToken(username string) string {
	token := uuid.New().String()

	s.mut.Lock()
	defer s.mut.Unlock()

	s.store[token] = &sessionData{username: username, updated: s.TimeProvider.Now()}

	return token
}

// InvalidateToken deletes any data associated with the given token. If the token is unknown
// it's a no-op.
func (s *MemoryStore) InvalidateToken(token string) {
	s.mut.Lock()
	defer s.mut.Unlock()

	delete(s.store, token)
}

func (s *MemoryStore) sessionDataForToken(token string) (*sessionData, error) {
	data, ok := s.store[token]

	if !ok {
		return nil, ErrInvalidToken
	}

	if s.TimeProvider.Since(data.updated) > s.TokenTimeout {
		delete(s.store, token)
		return nil, ErrInvalidToken
	}

	return data, nil
}

// IsTokenValid returns whether the given token is known and valid. Both unknown and timed out
// tokens are considered invalid.
func (s *MemoryStore) IsTokenValid(token string) bool {
	s.mut.Lock()
	defer s.mut.Unlock()

	_, err := s.sessionDataForToken(token)

	return err == nil
}

// GetUsernameForToken returns the associated username for a given token or an ErrInvalidToken
// error if the token is invalid.
func (s *MemoryStore) GetUsernameForToken(token string) (string, error) {
	s.mut.Lock()
	defer s.mut.Unlock()

	data, err := s.sessionDataForToken(token)

	if err != nil {
		return "", fmt.Errorf("could not get username for token: %w", err)
	}

	return data.username, nil
}

// ResetTokenTimeout resets the validity timeout of a given token, making it valid from the
// current time until TokenTimeout amount of time passed. Returns an ErrInvalidToken error if
// the token is invalid.
func (s *MemoryStore) ResetTokenTimeout(token string) error {
	s.mut.Lock()
	defer s.mut.Unlock()

	data, err := s.sessionDataForToken(token)

	if err != nil {
		return fmt.Errorf("could not get reset token timeout: %w", err)
	}

	data.updated = s.TimeProvider.Now()

	return nil
}

// New creates a new MemoryStore with a 60 minute default token timeout
func New() *MemoryStore {
	return &MemoryStore{
		store:        make(map[string]*sessionData),
		TimeProvider: realTimeProvider{},
		TokenTimeout: 60 * time.Minute,
	}
}
