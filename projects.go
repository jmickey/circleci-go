package circleci // import mickey.dev/go/circleci-go

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// ProjectService handles communication with the project
// related methods of the CircleCI API.
//
// https://circleci.com/docs/api/#projects
type ProjectService struct {
	client *Client
}

// Project type models the returned values from the
// /projects endpoint of the CircleCI API.
type Project struct {
	Name      string `json:"reponame"`
	Username  string `json:"username"`
	Following bool   `json:"following"`
	URL       string `json:"vcs_url"`
}

// List returns a slice containing all projects followed by the authenticated user.
func (p *ProjectService) List(ctx context.Context) ([]Project, error) {
	url := p.client.buildRequestURL(ProjectsEndpoint)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("Error retreiving projects from: %v: %w", url, err)
	}

	resp, err := p.client.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Unable to complete API request to %v: %w", url, err)
	}
	defer resp.Body.Close()

	var projects []Project
	body, _ := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(body, &projects)
	if err != nil {
		return nil, fmt.Errorf("Error decoding response body: '%v': %w", string(body), err)
	}

	return projects, nil
}

func (p *ProjectService) Get(ctx context.Context, proj string, username string) (*Project, error) {
	url := p.client.buildRequestURL(fmt.Sprintf("/project/github/%s/%s", username, proj))
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("Error retreiving project %s/%s from: %v: %w", username, proj, url, err)
	}

	resp, err := p.client.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Unable to complete API request to %v: %w", url, err)
	}
	defer resp.Body.Close()

	var project Project
	body, _ := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(body, &project)
	if err != nil {
		return nil, fmt.Errorf("Error decoding response body: '%v': %w", string(body), err)
	}

	return &project, nil
}
