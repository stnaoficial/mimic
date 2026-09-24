package cmd

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

const (
	DefaultConfigDirectoryPath          = ".mimic"
	DefaultConfigTemplatesDirectoryName = "templates"
	DefaultConfigTemplatesDirectoryPath = ".mimic/templates"
	DefaultConfigFilePath               = ".mimic/config"
)

var DefaultConfigSettings = map[string]string{
	"update.github.user.username":          "stnaoficial",
	"update.github.user.access-token":      "",
	"update.github.repository.name":        "mimic",
	"update.github.repository.branch.name": "main",

	"remote.github.user.username":          "stnaoficial",
	"remote.github.user.access-token":      "",
	"remote.github.repository.name":        "mimic-template",
	"remote.github.repository.branch.name": "main",
}

type Config struct {
	fileName string
	dirName  string

	settings map[string]string
}

func NewConfig(dirName string) *Config {
	fileName := filepath.Join(dirName, "config")

	var data = []byte{}

	if fileData, err := os.ReadFile(fileName); err == nil {
		data = fileData
	}

	settings := make(map[string]string)

	for line := range strings.SplitSeq(string(data), "\n") {
		if len(line) == 0 {
			continue
		}

		parts := strings.Split(line, ":")

		if len(parts) < 1 || len(parts) > 2 {
			continue
		}

		value := ""

		if len(parts) == 2 {
			value = strings.TrimSpace(parts[1])
		}

		key := strings.TrimSpace(parts[0])

		settings[key] = value
	}

	return &Config{
		fileName: fileName,
		dirName:  dirName,
		settings: settings,
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

func (c *Config) RegisterDefaultSettings() {
	for key, value := range DefaultConfigSettings {
		if _, ok := c.settings[key]; !ok {
			c.settings[key] = value
		}
	}
}

func (c *Config) Bytes() []byte {
	keys := make([]string, 0, len(c.settings))

	for key := range c.settings {
		keys = append(keys, key)
	}

	sort.Strings(keys)

	var buffer bytes.Buffer

	previousPrefix := ""

	for _, key := range keys {
		prefix := strings.SplitN(key, ".", 2)[0]

		if previousPrefix != "" && prefix != previousPrefix {
			buffer.WriteRune('\n')
		}

		buffer.WriteString(fmt.Sprintf("%s: %s\n", key, c.settings[key]))

		previousPrefix = prefix
	}

	buffer.WriteRune('\n')

	return buffer.Bytes()
}

func (c *Config) Require(key string) string {
	return c.settings[key]
}

func (c *Config) Get(key string) (string, bool) {
	value, ok := c.settings[key]
	return value, ok
}

func (c *Config) GetOrDefault(key string, defaultValue string) string {
	if value, ok := c.Get(key); !ok {
		return defaultValue
	} else {
		return value
	}
}

func (c *Config) GetAll(prefix string) map[string]string {
	settings := make(map[string]string)

	for key, value := range c.settings {
		if strings.HasPrefix(key, prefix) {
			settings[key] = value
		}
	}

	return settings
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
