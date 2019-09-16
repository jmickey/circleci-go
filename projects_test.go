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
		rw.Write([]byte(`[
			{
				"reponame": "somerepo",
				"following": true
			}
		]`))
	}))

	c, err := circleci.NewClient("someapitoken", srv.URL, circleci.SetBaseHTTPClient(srv.Client()))
	require.Nil(t, err, "circleci.NewClient() error should be nil, received: %v", err)

	projects, err := c.Projects.List(context.Background())
	require.Nil(t, err, "Received error when listing projects: %v", err)
	assert.NotNil(t, projects, "Projects is nil. Should contain testrepo item.")
	assert.Equal(t, "somerepo", projects[0].Name, "Project name did not equal 'somerepo'")
}

func TestGetProject(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		rw.WriteHeader(200)
		rw.Write([]byte(`[{
			"reponame": "testRepo",
			"following": true,
			"vcs_url": "https://github.com/testuser/testRepo",
			"username": "testuser"
		}]`))
	}))

	c, err := circleci.NewClient("someapitoken", srv.URL, circleci.SetBaseHTTPClient(srv.Client()))
	require.Nil(t, err, "circleci.NewClient() error should be nil, received: %v", err)

	project, err := c.Projects.Get(context.Background(), "testRepo", "testuser")
	require.Nil(t, err, "c.Projects.Get() returned an error, expected nil: %v", err)
	assert.Equal(t, "testRepo", project.Name, "Project.Name returned incorrect value. Expected: 'testRepo', got: %v", project.Name)
}
