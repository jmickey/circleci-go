package circleci // import "mickey.dev/go/circleci-go"

import (
	"context"
	"fmt"
)

// ProjectService handles communication with the project related methods of the
// CircleCI API. https://circleci.com/docs/api/#projects
type ProjectService struct {
	client *Client
}

// Project type models the returned values from the /projects endpoint of the
// CircleCI API.
type Project struct {
	Name     string `json:"reponame"`
	Username string `json:"username"`
	Followed bool   `json:"followed"`
	URL      string `json:"vcs_url"`
}

// List returns a slice containing all projects followed by the authenticated user.
func (p *ProjectService) List(ctx context.Context) ([]*Project, error) {
	urlPath := "projects"
	req, err := p.client.newRequestWithContext(ctx, "GET", urlPath, nil, nil)
	if err != nil {
		url := p.client.BaseURL.String() + urlPath
		return nil, fmt.Errorf("Error building request for: %v: %w", url, err)
	}

	var projects []*Project
	err = p.client.do(req, &projects)
	if err != nil {
		return nil, fmt.Errorf("Error completing API request to %v: %w", req.URL.String(), err)
	}

	return projects, nil
}

// Get retreives a project based on the project name and owner (GitHub Username).
// The authenticated user must be "following" the project in order to retrieve it.
func (p *ProjectService) Get(ctx context.Context, proj string, username string) (*Project, error) {
	projects, err := p.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("Could not retrieve projects, error: %w", err)
	}

	for _, project := range projects {
		if project.Name == proj && project.Username == username {
			return project, nil
		}
	}

	return nil, fmt.Errorf("Could not find project %s for user %s. Check you're following the project", proj, username)
}

func (p *ProjectService) Follow(ctx context.Context, project string, username string) error {
	urlPath := fmt.Sprintf("project/%s/%s/%s/follow", p.client.VCSType, username, project)
	req, err := p.client.newRequestWithContext(ctx, "POST", urlPath, nil, nil)
	if err != nil {
		url := p.client.BaseURL.String() + urlPath
		return fmt.Errorf("Error building request for: %v: %w", url, err)
	}

	err = p.client.do(req, nil)
	if err != nil {
		return fmt.Errorf("Error completing API request to %v: %w", req.URL.String(), err)
	}

	return nil
}

func (p *ProjectService) Unfollow(ctx context.Context, project string, username string) error {
	urlPath := fmt.Sprintf("project/%s/%s/%s/follow", p.client.VCSType, username, project)
	req, err := p.client.newRequestWithContext(ctx, "DELETE", urlPath, nil, nil)
	if err != nil {
		url := p.client.BaseURL.String() + urlPath
		return fmt.Errorf("Error building request for: %v: %w", url, err)
	}

	err = p.client.do(req, nil)
	if err != nil {
		return fmt.Errorf("Error completing API request to %v: %w", req.URL.String(), err)
	}

	return nil
}

// Enable the project in CircleCI, this will generate and add an SSH to the repo for code
// checkout. The authenticated user must have "admin" permissions on the rpeo.
func (p *ProjectService) Enable(ctx context.Context, project string, username string) error {
	urlPath := fmt.Sprintf("project/%s/%s/%s/enable", p.client.VCSType, username, project)
	req, err := p.client.newRequestWithContext(ctx, "POST", urlPath, nil, nil)
	if err != nil {
		url := p.client.BaseURL.String() + urlPath
		return fmt.Errorf("Error building request for: %v: %w", url, err)
	}

	err = p.client.do(req, nil)
	if err != nil {
		return fmt.Errorf("Error completing API request to %v: %w", req.URL.String(), err)
	}

	return nil
}

// Disable will remove the CircleCI deploy key from the repo. The authenicated user must
// have "admin" permissions on the repo.
func (p *ProjectService) Disable(ctx context.Context, project, username string) error {
	urlPath := fmt.Sprintf("project/%s/%s/%s/disable", p.client.VCSType, username, project)
	req, err := p.client.newRequestWithContext(ctx, "POST", urlPath, nil, nil)
	if err != nil {
		url := p.client.BaseURL.String() + urlPath
		return fmt.Errorf("Error building request for: %v: %w", url, err)
	}

	err = p.client.do(req, nil)
	if err != nil {
		return fmt.Errorf("Error completing API request to %v: %w", req.URL.String(), err)
	}

	return nil
}

// EnableAndFollow is a helper funtion that will enable the project, and then follow
// the enabled project.
func (p *ProjectService) EnableAndFollow(ctx context.Context, project, username string) error {
	err := p.Enable(ctx, project, username)
	if err != nil {
		return err
	}
	err = p.Follow(ctx, project, username)
	if err != nil {
		return err
	}
	return nil
}
