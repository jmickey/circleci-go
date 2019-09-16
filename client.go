package circleci // import mickey.dev/go/circleci-go

import (
	"fmt"
	"net/http"
	"net/url"
	"path"
)

var (
	ProjectsEndpoint = "/projects"
)

type Client struct {
	APIKey     string
	ServerURL  string
	httpClient *http.Client

	Projects *ProjectService
}

type ClientOption func(*Client) error

type ReqOption func(*http.Request) error

func NewClient(apiKey, serverURL string, opts ...ClientOption) (*Client, error) {
	apiPath := "/api/v1.1/"
	u, err := url.Parse(serverURL)
	if err != nil {
		return nil, fmt.Errorf("Couldn't parse URL: %v. Error: %w", serverURL, err)
	}

	u.Path = path.Join(u.Path, apiPath)

	c := &Client{
		APIKey:     apiKey,
		ServerURL:  serverURL,
		httpClient: &http.Client{},
	}

	for _, opt := range opts {
		err := opt(c)
		if err != nil {
			return nil, fmt.Errorf("Unable to configure client with options func: %w", err)
		}
	}

	c.Projects = &ProjectService{client: c}

	return c, nil
}

func (c *Client) buildRequestURL(endpoint string) string {
	return fmt.Sprintf("%s/%s?circle-token=%s", c.ServerURL, endpoint, c.APIKey)
}

func SetBaseHTTPClient(client *http.Client) func(*Client) error {
	return func(c *Client) error {
		c.httpClient = client
		return nil
	}
}

func SetHeaders(headers map[string]string) func(*http.Request) {
	return func(req *http.Request) {
		for k, v := range headers {
			req.Header.Set(k, v)
		}
	}
}

func (c *Client) GetHTTPClient() *http.Client {
	return c.httpClient
}
