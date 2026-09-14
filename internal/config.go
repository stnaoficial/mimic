package internal

import (
	"os"
	"path/filepath"
	"strings"
)

type Config struct {
	fileName string
	dirName  string
	entries  map[string]string
}

func NewConfig(dirName string) *Config {
	fileName := filepath.Join(dirName, "config")

	var data []byte

	if fileData, err := os.ReadFile(fileName); err == nil {
		data = fileData
	}

	entries := make(map[string]string)

	for line := range strings.SplitSeq(string(data), "\n") {
		parts := strings.Split(line, ":")

		if len(parts) < 1 || len(parts) > 2 {
			continue
		}

		value := ""

		if len(parts) == 2 {
			value = strings.TrimSpace(parts[1])
		}

		key := strings.TrimSpace(parts[0])

		entries[key] = value
	}

	return &Config{
		fileName: fileName,
		dirName:  dirName,
		entries:  entries,
	}
}

func NewLocalConfig() (*Config, error) {
	cwd, err := os.Getwd()

	if err != nil {
		return nil, err
	}

	homeDir, err := os.UserHomeDir()

	if err != nil {
		return nil, err
	}

	localConfigPath := filepath.Join(cwd, DefaultConfigDirectoryPath)
	globalConfigPath := filepath.Join(homeDir, DefaultConfigDirectoryPath)

	if _, err := os.ReadDir(localConfigPath); err == nil {
		return NewConfig(localConfigPath), nil
	}

	return NewConfig(globalConfigPath), nil
}

func (c *Config) GetOrDefault(key string, defaultValue string) string {
	if value, ok := c.entries[key]; ok {
		return value
	} else {
		return defaultValue
	}
}

func (c *Config) FileName() string {
	return c.fileName
}

func (c *Config) DirName() string {
	return c.dirName
}

func (c *Config) TemplatesDirectoryPath() string {
	return filepath.Join(c.dirName, DefaultConfigTemplatesDirectoryName)
}
