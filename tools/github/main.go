package github

import "net/url"

var PublicUrl url.URL = url.URL{Scheme: "https", Host: "github.com"}
var RawApiUrl url.URL = url.URL{Scheme: "https", Host: "raw.githubusercontent.com"}
var ApiUrl url.URL = url.URL{Scheme: "https", Host: "api.github.com"}
