package cache

import (
	"bytes"
	"encoding/gob"
	"errors"
	"fmt"
	"mimic/internal/util"
	"os"
	"path/filepath"
	"time"
)

type Storage struct {
	Name string
	Path string
}

type Payload struct {
	Value   []byte
	Expires time.Time
}

func NewStorage(name string) *Storage {
	return &Storage{
		Name: name,
		Path: filepath.Join(util.DefaultTempDir(), name),
	}
}

func (c *Storage) buildFileName(key *Key) string {
	return filepath.Join(c.Path, fmt.Sprintf("%x", key.Hash))
}

func (c *Storage) Backup(key *Key, value []byte, duration time.Duration) error {
	fileName := c.buildFileName(key)

	if err := os.MkdirAll(filepath.Dir(fileName), 0700); err != nil {
		return err
	}

	payload := Payload{
		Value:   value,
		Expires: time.Now().Add(duration),
	}

	var buf bytes.Buffer

	if err := gob.NewEncoder(&buf).Encode(payload); err != nil {
		return err
	}

	return os.WriteFile(fileName, buf.Bytes(), 0600)
}

func (c *Storage) Restore(key *Key) ([]byte, error) {
	value, err := os.ReadFile(c.buildFileName(key))

	if err != nil {
		return nil, err
	}

	var payload Payload

	if err := gob.NewDecoder(bytes.NewReader(value)).Decode(&payload); err != nil {
		return nil, err
	}

	if time.Now().After(payload.Expires) {
		return nil, errors.New("temporary Cache payload has expired")
	}

	return payload.Value, nil
}
