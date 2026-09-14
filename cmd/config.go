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

	Settings map[string]string
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

	config := &Config{
		fileName: fileName,
		dirName:  dirName,
		Settings: settings,
	}

	for key, value := range DefaultConfigSettings {
		if _, ok := config.Settings[key]; !ok {
			config.Settings[key] = value
		}
	}

	return config
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

func (c *Config) Bytes() []byte {
	keys := make([]string, 0, len(c.Settings))

	for key := range c.Settings {
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

		buffer.WriteString(fmt.Sprintf("%s: %s\n", key, c.Settings[key]))

		previousPrefix = prefix
	}

	buffer.WriteRune('\n')

	return buffer.Bytes()
}

func (c *Config) Get(key string) string {
	return c.Settings[key]
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
