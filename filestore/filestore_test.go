package filestore

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFilestoreHandlesEscapingPath(t *testing.T) {
	fs := NewLocal(os.TempDir())
	_, _, err := fs.GetUserFile(context.Background(), "test", "..")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "path may not contain ..")
}
