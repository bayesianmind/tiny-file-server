package userstore

import (
	"bytes"
	"context"
	"fmt"
	"sync"

	"github.com/bayesianmind/demo-file-server/auth"
)

type inmem struct {
	users map[string][]byte
	usermu sync.RWMutex
}

func NewInMemory() Interface {
	return &inmem{
		users: make(map[string][]byte),
	}
}

func (m *inmem) RegisterUser(ctx context.Context, user, password string) error {
	m.usermu.Lock()
	defer m.usermu.Unlock()

	hashed := hashPw(password)

	if dbPw, found := m.users[user]; found && !bytes.Equal(dbPw, hashed) {
		// idempotent only if the password is the same
		return fmt.Errorf("username already registered")
	}

	m.users[user] = hashed
	return nil
}

func (m *inmem) Login(ctx context.Context, user, password string) (*LoginResponse, error) {
	m.usermu.RLock()
	defer m.usermu.RUnlock()

	hashed := hashPw(password)

	if dbPw, found := m.users[user]; !found || !bytes.Equal(dbPw, hashed) {
		return nil, fmt.Errorf("invalid login")
	}

	token, err := auth.TokenForUser(user)
	resp := &LoginResponse{
		SessionToken: token,
	}

	return resp, err
}
