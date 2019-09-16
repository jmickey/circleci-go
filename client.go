package circleci // import "mickey.dev/go/circleci-go"

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
)

var (
	ProjectsEndpoint = "projects"
)

type Client struct {

	// APIKey stores the CircleCI API Token. Generate a CircleCI token by
	// navigating to https://circleci.com/account/api and generating a new token.
	//
	// The API token is required in order to communicate with the CircleCI
	// API server. The authenticated user should be "Following" the projects
	// They want to make changes to.
	//
	// New projects can also be followed by using the client.Projects.Follow() API call.
	APIKey string

	// ServerURL should be the CircleCI server URL to make API calls against.
	//
	// For CircleCI SaaS use https://circleci.com. For CircleCI Server
	// (previously Enterprise) use your internal CircleCI URL. If your
	// CircleCI Server requires a custom port that isn't 443, then include
	// this in the URL. E.g. https://circle.company.com:8080
	ServerURL  string
	httpClient *http.Client

	// Projects represents the CircleCI project resources. This will be
	// instantiated by default with a *circleci.ProjectService object when
	// circleci.NewClient() is called.
	Projects *ProjectService
}

type ClientOption func(*Client) error

type ReqOption func(*http.Request) error

func NewClient(apiKey, server string, opts ...ClientOption) (*Client, error) {
	apiPath := "/api/v1.1/"
	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, fmt.Errorf("Couldn't parse URL: %v. Error: %w", server, err)
	}

	serverURL.Path = path.Join(serverURL.Path, apiPath)

	c := &Client{
		APIKey:     apiKey,
		ServerURL:  serverURL.String(),
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

func (c *Client) newCircleRequestWithContext(ctx context.Context, verb string, url string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequestWithContext(ctx, verb, url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Accept", "applications/json")
	return req, nil
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
