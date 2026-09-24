package temp

import (
	"crypto/sha256"
)

type Key struct {
	value string
	Hash  [32]byte
}

func NewKey(value string) *Key {
	return &Key{
		value: value,
		Hash:  sha256.Sum256([]byte(value)),
	}
}
