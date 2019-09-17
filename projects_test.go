// +build !integration

package circleci_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"mickey.dev/go/circleci-go"
)

func TestListProjects(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		rw.WriteHeader(200)
		rw.Write([]byte(listProjectsResponse))
	}))
	defer srv.Close()

	c, err := circleci.NewClient("someapitoken", circleci.WithBaseServerURL(srv.URL), circleci.WithBaseHTTPClient(srv.Client()))
	require.Nil(t, err, "circleci.NewClient() error should be nil, received: %v", err)

	projects, err := c.Projects.List(context.Background())
	require.Nil(t, err, "Received error when listing projects: %v", err)
	assert.NotNil(t, projects, "Projects is nil. Should contain testrepo item.")
	assert.Equal(t, "reponame", projects[0].Name, "Project.Name returned incorrect value. Expected: 'reponame', got: %v", projects[0].Name)
}

func TestGetProject(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		rw.WriteHeader(200)
		rw.Write([]byte(listProjectsResponse))
	}))
	defer srv.Close()

	c, err := circleci.NewClient("someapitoken", circleci.WithBaseServerURL(srv.URL), circleci.WithBaseHTTPClient(srv.Client()))
	require.Nil(t, err, "circleci.NewClient() error should be nil, received: %v", err)

	project, err := c.Projects.Get(context.Background(), "reponame", "username")
	require.Nil(t, err, "c.Projects.Get() returned an error, expected nil: %v", err)
	assert.Equal(t, "reponame", project.Name, "Project.Name returned incorrect value. Expected: 'reponame', got: %v", project.Name)
}

func TestFollowProject(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		result := assert.Equal(t, "POST", r.Method, "Unexpected reqest method. Expected: %s. Received: %s", "POST", r.Method)
		if result {
			rw.WriteHeader(200)
		} else {
			rw.WriteHeader(500)
		}
	}))
	defer srv.Close()

	c, err := circleci.NewClient("SomeAPIToken", circleci.WithBaseServerURL(srv.URL), circleci.WithBaseHTTPClient(srv.Client()))
	require.Nil(t, err, "circleci.NewClient() error should be nil, received: %v", err)

	err = c.Projects.Follow(context.Background(), "reponame", "username")
	require.Nil(t, err, "c.Projects.Follow returned an error, expected nil: %v", err)
}

var (
	listProjectsResponse = `[{
		"reponame": "reponame",
		"following": true,
		"vcs_url": "https://github.com/testuser/testRepo",
		"username": "username"
	}]`
)
