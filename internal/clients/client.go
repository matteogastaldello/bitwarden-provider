package bwclient

import (
	"context"
	"fmt"
	"net/http"
	"path"

	"github.com/lucasepe/httplib"
)

const (
	ApiVersionKey  = "api-version"
	ApiVersionVal  = "7.0"
	ApiPreviewFlag = "-preview"
	UserAgent      = "krateo/azuredevops-provider"
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

func NewClient() *Client {
	return &Client{
		httpClient: httplib.NewClient(),
		uriMap: map[URIKey]string{
			Default: "http://host.docker.internal:8087",
		},
	}
}

func Unlock(ctx context.Context, c *Client, password string) (*Response, error) {
	uri, err := httplib.NewURLBuilder(httplib.URLBuilderOptions{
		BaseURL: c.BaseURL(Default),
		Path:    path.Join("unlock"),
	}).Build()

	if err != nil {
		return nil, err
	}
	values := struct {
		Password string `json:"password"`
	}{
		Password: password,
	}
	req, err := httplib.Post(uri.String(), httplib.ToJSON(values))
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json")
	req = req.WithContext(ctx)

	apiErr := &APIError{}

	val := &Response{}

	err = httplib.Fire(c.HTTPClient(), req, httplib.FireOptions{
		ResponseHandler: httplib.FromJSON(val),
		Validators: []httplib.HandleResponseFunc{
			httplib.ErrorJSON(apiErr, http.StatusOK),
		},
	})
	return val, err
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
