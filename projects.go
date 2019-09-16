package circleci // import mickey.dev/go/circleci-go

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
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
	req, err := p.client.newCircleRequestWithContext(ctx, "GET", url, nil)
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

// Get retreives a project based on the project name and owner (GitHub Username).
// The authenticated user must be "following" the project in order to retreive it.
func (p *ProjectService) Get(ctx context.Context, proj string, username string) (*Project, error) {
	projects, err := p.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("Could not retrieve projects, error: %w", err)
	}

	for _, project := range projects {
		if project.Name == proj && project.Username == username {
			return &project, nil
		}
	}

	return nil, fmt.Errorf("Could not find project %s for user %s. Check you're following the project", proj, username)
}
