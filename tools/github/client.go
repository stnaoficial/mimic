package github

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"mimic/tools/temp"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"
	"time"
)

const (
	ApiTreeEntryTypeTree = "tree"
	ApiTreeEntryTypeBlob = "blob"

	ApiBlobEntryEncodingBase64 = "base64"
)

type ApiTreeEntry struct {
	Path string `json:"path"`
	Mode string `json:"mode"`
	Type string `json:"type"`
	Sha  string `json:"sha"`
	Size int64  `json:"size"`
	Url  string `json:"url"`
}

type ApiTreeResponse struct {
	Sha       string         `json:"sha"`
	Tree      []ApiTreeEntry `json:"tree"`
	Truncated bool           `json:"truncated"`
}

type ApiBlobEntry struct {
	Sha      string `json:"sha"`
	NodeId   string `json:"node_id"`
	Size     int64  `json:"size"`
	Url      string `json:"url"`
	Content  string `json:"content"`
	Encoding string `json:"encoding"`
}

func (e *ApiBlobEntry) DecodeContent() ([]byte, error) {
	if e.Encoding != ApiBlobEntryEncodingBase64 {
		return nil, fmt.Errorf("Unexpected encoding type %s", e.Encoding)
	}

	content := strings.ReplaceAll(e.Content, "\n", "")

	decoded, err := base64.StdEncoding.DecodeString(content)

	if err != nil {
		return nil, err
	}

	return decoded, nil
}

type ApiEntryLink struct {
	Self string `json:"self"`
	Git  string `json:"git"`
	Html string `json:"html"`
}

type ApiEntry struct {
	Name        string         `json:"name"`
	Path        string         `json:"path"`
	Sha         string         `json:"sha"`
	Size        int64          `json:"size"`
	Url         string         `json:"url"`
	HtmlUrl     string         `json:"html_url"`
	GitUrl      string         `json:"git_url"`
	DownloadUrl string         `json:"download_url"`
	Type        string         `json:"type"`
	Links       []ApiEntryLink `json:"links"`
}

type ApiResponse struct {
	SHA       string     `json:"sha"`
	Tree      []ApiEntry `json:"tree"`
	Truncated bool       `json:"truncated"`
}

type ApiClient struct {
	PublicUrl url.URL
	RawApiUrl url.URL
	ApiUrl    url.URL

	tempStorage *temp.Storage

	username       string
	branchName     string
	repositoryName string
	accessToken    string

	Cache bool
}

func NewApiClient(username string, repositoryName string, branchName string, accessToken string) *ApiClient {
	return &ApiClient{
		PublicUrl: *PublicUrl.JoinPath(username, repositoryName),
		RawApiUrl: *RawApiUrl.JoinPath(username, repositoryName),
		ApiUrl:    *ApiUrl.JoinPath("repos", username, repositoryName),

		tempStorage: temp.NewStorage("github"),

		username:       username,
		repositoryName: repositoryName,
		branchName:     branchName,
		accessToken:    accessToken,

		Cache: true,
	}
}

func (g *ApiClient) Username() string {
	return g.username
}

func (g *ApiClient) BranchName() string {
	return g.branchName
}

func (g *ApiClient) RepositoryName() string {
	return g.repositoryName
}

func (g *ApiClient) ParsePublicEntryUrl(path string) *url.URL {
	return g.PublicUrl.JoinPath("tree", g.branchName, path)
}

func (g *ApiClient) ParseRawApiEntryUrl(path string) *url.URL {
	return g.RawApiUrl.JoinPath(g.branchName, path)
}

func (g *ApiClient) ParseRemoteEntry(entry ApiTreeEntry) *RemoteEntry {
	info := &RemoteFileInfo{
		name:    path.Base(entry.Path),
		size:    entry.Size,
		modTime: time.Time{},
		mode:    parseMode(entry.Mode),
		isDir:   entry.Type == ApiTreeEntryTypeTree,
	}

	return &RemoteEntry{
		Path: entry.Path,
		Info: info,
	}
}

func (g *ApiClient) Walk(target string) ([]ApiTreeEntry, error) {
	tree, err := g.FetchApiTree()

	if err != nil {
		return nil, err
	}

	var entries []ApiTreeEntry

	for _, entry := range tree.Tree {
		if !isPathInside(target, entry.Path) {
			continue
		}

		entries = append(entries, entry)
	}

	if tree.Truncated {
		return nil, fmt.Errorf("Tree is truncated")
	}

	if len(entries) == 0 {
		return nil, fmt.Errorf("Entry not found")
	}

	hasChildren := make(map[string]bool)

	for _, entry := range entries {
		parent := path.Dir(entry.Path)

		if parent == "." {
			parent = ""
		}

		for parent != target && parent != "" && parent != "/" {
			hasChildren[parent] = true
			parent = path.Dir(parent)

			if parent == "." {
				parent = ""
			}
		}

		hasChildren[parent] = true
	}

	var result []ApiTreeEntry

	for _, entry := range entries {
		if entry.Type == ApiTreeEntryTypeTree && hasChildren[entry.Path] {
			continue
		}

		result = append(result, entry)
	}

	return result, nil
}

func (g *ApiClient) FetchApi(url *url.URL) (*http.Response, error) {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	req, err := http.NewRequest("GET", url.String(), nil)

	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("X-GitHub-Api-Version", "2026-03-10")

	if len(g.accessToken) != 0 {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", g.accessToken))
	}

	return client.Do(req)
}

func (g *ApiClient) FetchApiTree() (*ApiTreeResponse, error) {
	url := g.ApiUrl.JoinPath("git", "trees", g.branchName)

	query := url.Query()
	query.Set("recursive", "1")

	url.RawQuery = query.Encode()

	tempKey := temp.NewKey(url.String())

	if value, err := g.tempStorage.Restore(tempKey); err == nil && g.Cache {
		return g.parseTreeResponse(bytes.NewReader(value))
	}

	res, err := g.FetchApi(url)

	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Invalid status code %s", res.Status)
	}

	body, err := io.ReadAll(res.Body)

	if err != nil {
		return nil, err
	}

	if err := g.tempStorage.Backup(tempKey, body, time.Minute); err != nil {
		return nil, err
	}

	return g.parseTreeResponse(bytes.NewReader(body))
}

func (g *ApiClient) FetchApiContent(path string) ([]ApiEntry, error) {
	url := g.ApiUrl.JoinPath("contents", path)

	tempKey := temp.NewKey(url.String())

	if value, err := g.tempStorage.Restore(tempKey); err == nil && g.Cache {
		return g.parseContentResponse(bytes.NewReader(value))
	}

	res, err := g.FetchApi(url)

	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Invalid status code %s", res.Status)
	}

	body, err := io.ReadAll(res.Body)

	if err != nil {
		return nil, err
	}

	if err := g.tempStorage.Backup(tempKey, body, time.Minute); err != nil {
		return nil, err
	}

	return g.parseContentResponse(bytes.NewReader(body))
}

func (g *ApiClient) FetchApiBlob(sha string) (*ApiBlobEntry, error) {
	url := g.ApiUrl.JoinPath("git", "blobs", sha)

	tempKey := temp.NewKey(url.String())

	if value, err := g.tempStorage.Restore(tempKey); err == nil && g.Cache {
		return g.parseBlobResponse(bytes.NewReader(value))
	}

	res, err := g.FetchApi(url)

	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Invalid status code %s\n", res.Status)
	}

	body, err := io.ReadAll(res.Body)

	if err != nil {
		return nil, err
	}

	if err := g.tempStorage.Backup(tempKey, body, time.Minute); err != nil {
		return nil, err
	}

	return g.parseBlobResponse(bytes.NewReader(body))
}

func (g *ApiClient) FetchRawApiContent(path string) ([]byte, error) {
	url := g.RawApiUrl.JoinPath(g.branchName, path)

	tempKey := temp.NewKey(url.String())

	if value, err := g.tempStorage.Restore(tempKey); err == nil && g.Cache {
		return value, nil
	}

	res, err := g.FetchApi(url)

	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Invalid status code %s\n", res.Status)
	}

	body, err := io.ReadAll(res.Body)

	if err != nil {
		return nil, err
	}

	if err := g.tempStorage.Backup(tempKey, body, time.Minute); err != nil {
		return nil, err
	}

	return body, nil
}

func (g *ApiClient) parseTreeResponse(value io.Reader) (*ApiTreeResponse, error) {
	var tree ApiTreeResponse

	if err := json.NewDecoder(value).Decode(&tree); err != nil {
		return nil, err
	}

	return &tree, nil
}

func (g *ApiClient) parseBlobResponse(value io.Reader) (*ApiBlobEntry, error) {
	var res ApiBlobEntry

	if err := json.NewDecoder(value).Decode(&res); err != nil {
		return nil, err
	}

	return &res, nil
}

func (g *ApiClient) parseContentResponse(value io.Reader) ([]ApiEntry, error) {
	var res []ApiEntry

	if err := json.NewDecoder(value).Decode(&res); err != nil {
		return nil, err
	}

	return res, nil
}

func parseMode(mode string) os.FileMode {
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
	rootSegments := strings.Split(strings.Trim(root, "/"), "/")
	targetSegments := strings.Split(strings.Trim(target, "/"), "/")

	if len(rootSegments) > len(targetSegments) {
		return false
	}

	for index, rootSegment := range rootSegments {
		if rootSegment != targetSegments[index] {
			return false
		}
	}

	return true
}
