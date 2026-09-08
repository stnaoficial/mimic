package util

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"strings"
	"time"
)

type GitHubUtility struct {
	treeBaseUrl string
	rawBaseUrl  string
}

type GitHubRemoteEntry struct {
	Path string
	Info os.FileInfo
}

type GitHubRemoteEntryMap map[string]GitHubRemoteEntry

type GitHubRemoteFileInfo struct {
	name    string
	size    int64
	modTime time.Time
	mode    os.FileMode
	isDir   bool
}

func (f *GitHubRemoteFileInfo) Name() string       { return f.name }
func (f *GitHubRemoteFileInfo) Size() int64        { return f.size }
func (f *GitHubRemoteFileInfo) Mode() os.FileMode  { return f.mode }
func (f *GitHubRemoteFileInfo) ModTime() time.Time { return f.modTime }
func (f *GitHubRemoteFileInfo) IsDir() bool        { return f.isDir }
func (f *GitHubRemoteFileInfo) Sys() any           { return nil }

type GitHubTreeEntry struct {
	Path string `json:"path"`
	Mode string `json:"mode"`
	Type string `json:"type"`
	SHA  string `json:"sha"`
	Size int64  `json:"size"`
}

type GitHubTreeResponse struct {
	SHA       string            `json:"sha"`
	Tree      []GitHubTreeEntry `json:"tree"`
	Truncated bool              `json:"truncated"`
}

func NewGitHubUtility(userName string, repoName string) *GitHubUtility {
	return &GitHubUtility{
		treeBaseUrl: "https://api.github.com/repos/" + userName + "/" + repoName + "/git/trees",
		rawBaseUrl:  "https://raw.githubusercontent.com/" + userName + "/" + repoName + "/main",
	}
}

func (g *GitHubUtility) Walk(root string) ([]GitHubRemoteEntry, error) {
	tree, err := g.FetchTree()

	if err != nil {
		return nil, err
	}

	var entries []GitHubRemoteEntry

	for _, entry := range tree.Tree {
		if !isPathInside(root, entry.Path) {
			continue
		}

		info := &GitHubRemoteFileInfo{
			name:    path.Base(entry.Path),
			size:    entry.Size,
			modTime: time.Time{},
			mode:    githubMode(entry.Mode),
			isDir:   entry.Type == "tree",
		}

		entries = append(entries, GitHubRemoteEntry{
			Path: entry.Path,
			Info: info,
		})
	}

	if tree.Truncated {
		return nil, fmt.Errorf("Github repository tree is truncated")
	}

	return filterRemoteEntries(root, entries), nil
}

func (g *GitHubUtility) FetchTree() (*GitHubTreeResponse, error) {
	url := g.treeBaseUrl + "/main?recursive=1"

	res, err := http.Get(url)

	if err != nil {
		return nil, fmt.Errorf("Failed to fetch \"%s\"", url)
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Invalid status code %s", res.Status)
	}

	var tree GitHubTreeResponse

	if err := json.NewDecoder(res.Body).Decode(&tree); err != nil {
		return nil, fmt.Errorf("Failed to decode GitHub tree")
	}

	return &tree, nil
}

func (g *GitHubUtility) FetchRawContent(filePath string) ([]byte, error) {
	url := g.rawBaseUrl + "/" + filePath

	res, err := http.Get(url)

	if err != nil {
		return nil, fmt.Errorf("Failed to fetch \"%s\"\n", url)
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Invalid status code %s\n", res.Status)
	}

	body, err := io.ReadAll(res.Body)

	if err != nil {
		return nil, fmt.Errorf("Failed to read response body")
	}

	return body, nil
}

func githubMode(mode string) os.FileMode {
	switch mode {
	case "100755":
		return 0755
	case "100644":
		return 0644
	case "040000":
		return os.ModeDir | 0755
	default:
		return 0444
	}
}

func isPathInside(root string, target string) bool {
	root = strings.Trim(root, "/")
	target = strings.Trim(target, "/")

	if root == "" {
		return true
	}

	return root == target || strings.HasPrefix(target, root)
}

func filterRemoteEntries(root string, entries []GitHubRemoteEntry) []GitHubRemoteEntry {
	hasChildren := make(map[string]bool)

	for _, entry := range entries {
		parent := path.Dir(entry.Path)

		if parent == "." {
			parent = ""
		}

		for parent != root && parent != "" && parent != "/" {
			hasChildren[parent] = true
			parent = path.Dir(parent)

			if parent == "." {
				parent = ""
			}
		}

		hasChildren[parent] = true
	}

	var result []GitHubRemoteEntry

	for _, entry := range entries {
		if entry.Info.IsDir() && hasChildren[entry.Path] {
			continue
		}

		result = append(result, entry)
	}

	return result
}
