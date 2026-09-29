package github

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mimic/tools/cache"
	"net/http"
	"net/url"
	"path/filepath"
	"time"
)

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

type ApiEntryResponse struct {
	SHA       string     `json:"sha"`
	Tree      []ApiEntry `json:"tree"`
	Truncated bool       `json:"truncated"`
}

type ApiClient struct {
	PublicUrl url.URL
	RawApiUrl url.URL
	ApiUrl    url.URL

	cacheStorage *cache.Storage

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

		cacheStorage: cache.NewStorage(filepath.Join("github", "cache")),

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

func (g *ApiClient) FetchApiContents(path string) ([]ApiEntry, error) {
	url := g.ApiUrl.JoinPath("contents", path)

	cacheKey := cache.NewKey(url.String())

	if value, err := g.cacheStorage.Restore(cacheKey); err == nil && g.Cache {
		return g.parseApiEntryResponseList(bytes.NewReader(value))
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

	if err := g.cacheStorage.Backup(cacheKey, body, time.Minute); err != nil {
		return nil, err
	}

	return g.parseApiEntryResponseList(bytes.NewReader(body))
}

func (g *ApiClient) parseApiEntryResponseList(value io.Reader) ([]ApiEntry, error) {
	var res []ApiEntry

	if err := json.NewDecoder(value).Decode(&res); err != nil {
		return nil, err
	}

	return res, nil
}
