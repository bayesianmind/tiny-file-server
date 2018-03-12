package userstore

import (
	"context"
	"crypto/sha256"
)

type Interface interface {
	RegisterUser(ctx context.Context, user, password string) error
	Login(ctx context.Context, user, password string) (*LoginResponse, error)
}

type LoginResponse struct {
	// SessionToken is an opaque session token
	SessionToken string
}

// salt our userstore to make generalized precomputed sha256 reverses useless
var salt = []byte("sdkhfkajsdhfjkxcbuiqwoiurefkjshdfkasjhdfkjsahdf")

func hashPw(password string) []byte {
	h := sha256.New()
	h.Write(salt)
	h.Write([]byte(password))
	return h.Sum(nil)
}