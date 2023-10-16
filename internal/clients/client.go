package bwclient

import (
	"fmt"
	"net/http"

	"github.com/lucasepe/httplib"
)

type URIKey string

const (
	Default URIKey = "default"
)

type Client struct {
	httpClient *http.Client
	uriMap     map[URIKey]string
}

type Response struct {
	Success bool
	Data    map[string]interface{}
}

type ClientOptions struct {
	ApiUrl   string
	Password string
}

func NewClient(apiURL string) *Client {
	return &Client{
		httpClient: httplib.NewClient(),
		uriMap: map[URIKey]string{
			Default: apiURL,
		},
	}
}

func (c *Client) BaseURL(loc URIKey) string {
	val, ok := c.uriMap[loc]
	if !ok {
		return c.uriMap[Default]
	}
	return val
}

func (c *Client) HTTPClient() *http.Client {
	return c.httpClient
}

type APIError struct {
	Message string `json:"message"`
}

func (e *APIError) Error() string {
	return fmt.Sprintf("Bitwarden cli-API: %s", e.Message)
}
