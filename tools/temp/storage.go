package temp

import (
	"bytes"
	"encoding/gob"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

type Storage struct {
	name string
}

type StoragePayload struct {
	Value   []byte
	Expires time.Time
}

func NewStorage(name string) *Storage {
	return &Storage{name}
}

func (s *Storage) buildFileName(key *Key) string {
	return filepath.Join(os.TempDir(), "mimic", s.name, fmt.Sprintf("%x", key.Hash))
}

func (s *Storage) Backup(key *Key, value []byte, duration time.Duration) error {
	fileName := s.buildFileName(key)

	if err := os.MkdirAll(filepath.Dir(fileName), 0700); err != nil {
		return err
	}

	payload := StoragePayload{
		Value:   value,
		Expires: time.Now().Add(duration),
	}

	var buf bytes.Buffer

	if err := gob.NewEncoder(&buf).Encode(payload); err != nil {
		return err
	}

	return os.WriteFile(fileName, buf.Bytes(), 0600)
}

func (s *Storage) Restore(key *Key) ([]byte, error) {
	value, err := os.ReadFile(s.buildFileName(key))

	if err != nil {
		return nil, err
	}

	var payload StoragePayload

	if err := gob.NewDecoder(bytes.NewReader(value)).Decode(&payload); err != nil {
		return nil, err
	}

	if time.Now().After(payload.Expires) {
		return nil, errors.New("temporary storage payload has expired")
	}

	return payload.Value, nil
}
