package circleci // import "mickey.dev/go/circleci-go"

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

const (
	defaultBaseURL = "https://circleci.com/api/v1.1/"
	userAgent      = "circleci-go"
	acceptHeader   = "application/json"
	contentType    = "application/json"
)

type Client struct {
	httpClient *http.Client
	// APIKey stores the CircleCI API Token. Generate a CircleCI token by navigating
	// to https://circleci.com/account/api and generating a new token.
	//
	// The API token is required in order to communicate with the CircleCI API server.
	// The authenticated user should be "Following" the projects they want to make
	// changes to.
	//
	// New projects can also be followed by using the client.Projects.Follow() API.
	APIKey string

	// ServerURL should be the CircleCI server URL to make API calls against.
	//
	// For CircleCI SaaS use https://circleci.com/. For CircleCI Server (previously
	// Enterprise) use your internal CircleCI URL. If your CircleCI Server requires a
	// custom port that isn't 443, then include this in the URL.
	// E.g. https://circle.company.com:8080/
	BaseURL *url.URL

	// Projects represents the CircleCI project resources. This will be instantiated
	// by default with a *circleci.ProjectService object when circleci.NewClient() is
	// called.
	Projects *ProjectService
}

type APIError struct {
	StatusCode int    `json:"statusCode,omitempty"`
	Message    string `json:"message"`
}

func (e *APIError) Error() string {
	return fmt.Sprintf("API returned error: %d %s", e.StatusCode, e.Message)
}

type ClientOption func(*Client) error

func WithBaseHTTPClient(client *http.Client) func(*Client) error {
	return func(c *Client) error {
		c.httpClient = client
		return nil
	}
}

func WithBaseServerURL(u string) func(*Client) error {
	return func(c *Client) error {
		if !strings.HasSuffix(u, "/") {
			u = u + "/"
		}
		baseURL, err := url.Parse(u)
		if err != nil {
			return fmt.Errorf("Could not parse base URL: %s: %w", u, err)
		}

		baseURL.Path = "api/v1.1/"
		c.BaseURL = baseURL
		return nil
	}
}

func NewClient(apiKey string, opts ...ClientOption) (*Client, error) {
	baseURL, err := url.Parse(defaultBaseURL)
	if err != nil {
		return nil, fmt.Errorf("Couldn't parse URL: %v. Error: %w", baseURL.String(), err)
	}

	c := &Client{
		APIKey:     apiKey,
		BaseURL:    baseURL,
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

func (c *Client) newRequestWithContext(ctx context.Context, verb string, urlPath string, params url.Values, body interface{}) (*http.Request, error) {
	if params == nil {
		params = url.Values{}
	}
	params.Set("circle-token", c.APIKey)

	u := c.BaseURL.ResolveReference(&url.URL{Path: urlPath, RawQuery: params.Encode()})

	req, err := http.NewRequestWithContext(ctx, verb, u.String(), nil)
	if err != nil {
		return nil, err
	}

	if body != nil {
		req.Header.Set("Content-Type", contentType)
	}
	req.Header.Set("Accept", acceptHeader)
	req.Header.Set("User-Agent", userAgent)
	return req, nil
}

func (c *Client) GetHTTPClient() *http.Client {
	return c.httpClient
}

func (c *Client) do(req *http.Request, v interface{}) error {
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("Error reading response body: %w", err)
		}

		if len(body) > 0 {
			apiError := &APIError{}
			err = json.Unmarshal(body, &apiError)
			if err != nil {
				return err
			}
			apiError.StatusCode = resp.StatusCode
			return apiError
		}

		return &APIError{StatusCode: resp.StatusCode}
	}

	if v != nil {
		err = json.NewDecoder(resp.Body).Decode(&v)
		if err != nil {
			return err
		}
	}

	return nil
}
