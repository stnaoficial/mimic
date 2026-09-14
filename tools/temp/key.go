package temp

import (
	"crypto/sha256"
	"time"
)

type Key struct {
	value   string
	expires time.Duration
	Hash    [32]byte
}

func NewKey(value string, expires time.Duration) *Key {
	return &Key{
		value:   value,
		expires: expires,
		Hash:    sha256.Sum256([]byte(value + expires.String())),
	}
}
