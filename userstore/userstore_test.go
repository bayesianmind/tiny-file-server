package userstore

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHashPw(t *testing.T) {
	hBytes := hashPw("testpw")
	assert.Equal(t, "a4f3c0ddb4a7fbbc1b71f7c54d00d94e5a55823e1db927116bf7427ab81aab7f", hex.EncodeToString(hBytes))
	// ^ sample value from https://passwordsgenerator.net/sha256-hash-generator/
}