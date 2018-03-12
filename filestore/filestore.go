package filestore

import (
	"context"
	"io"
)

type Interface interface {
	PutUserFile(ctx context.Context, user, path string, file io.Reader, contentType string) error
	// GetUserFile retrieves the user's file. Caller is expected to close the resulting Reader.
	GetUserFile(ctx context.Context, user, path string) (io.ReadCloser, string, error)
	DeleteUserFile(ctx context.Context, user, path string) error
	ListFiles(ctx context.Context, user string) ([]string, error)
}
