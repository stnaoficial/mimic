package temp

import (
	"fmt"
	"os"
	"path/filepath"
)

type Storage struct {
	name string
}

func NewStorage(name string) *Storage {
	return &Storage{
		name: name,
	}
}

func (s *Storage) parseFileName(key *Key) string {
	return filepath.Join(os.TempDir(), "mimic", s.name, fmt.Sprintf("%x", key.Hash))
}

func (s *Storage) Backup(key *Key, value []byte) error {
	fileName := s.parseFileName(key)
	dirName := filepath.Dir(fileName)

	if err := os.MkdirAll(dirName, 0700); err != nil {
		return err
	}

	if err := os.WriteFile(fileName, value, 0600); err != nil {
		return err
	}

	return nil
}

func (s *Storage) Restore(key *Key) ([]byte, error) {
	value, err := os.ReadFile(s.parseFileName(key))

	if err != nil {
		return nil, err
	}

	return value, nil
}
