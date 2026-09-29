package github

import (
	"fmt"
	"mimic/internal/util"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type Mirror struct {
	username       string
	branchName     string
	repositoryName string

	tempDir string
}

func NewMirror(username string, repositoryName string, branchName string) *Mirror {
	return &Mirror{
		username:       username,
		repositoryName: repositoryName,
		branchName:     branchName,

		tempDir: filepath.Join(util.DefaultTempDir(), "github", "mirror"),
	}
}

func (m *Mirror) Mirror(paths ...string) ([]string, error) {
	if err := os.RemoveAll(m.tempDir); err != nil {
		return nil, err
	}

	if err := os.MkdirAll(m.tempDir, 0755); err != nil {
		return nil, err
	}

	ssh := fmt.Sprintf("git@github.com:%s/%s.git", m.username, m.repositoryName)

	cmd := exec.Command("git", "clone", "--depth=1", "--filter=blob:none", "--sparse", "--branch", m.branchName, ssh, m.tempDir)

	if err := cmd.Run(); err != nil {
		return nil, err
	}

	cmd = exec.Command("git", "sparse-checkout", "init", "--no-cone")
	cmd.Dir = m.tempDir

	if err := cmd.Run(); err != nil {
		return nil, err
	}

	patterns := make([]string, 0, len(paths))

	for _, path := range paths {
		path = filepath.ToSlash(filepath.Clean(path))

		patterns = append(patterns, "/"+path)
		patterns = append(patterns, "/"+path+"/**")
	}

	cmd = exec.Command("git", "sparse-checkout", "set", "--no-cone", "--stdin")
	cmd.Dir = m.tempDir
	cmd.Stdin = strings.NewReader(strings.Join(patterns, "\n") + "\n")

	if err := cmd.Run(); err != nil {
		return nil, err
	}

	results := make([]string, 0, len(paths))

	for _, path := range paths {
		results = append(results, filepath.Join(m.tempDir, filepath.FromSlash(filepath.Clean(path))))
	}

	return results, nil
}
