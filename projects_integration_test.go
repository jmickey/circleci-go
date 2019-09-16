// +build integration

package circleci_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"mickey.dev/go/circleci-go"
)

func TestListProjectsIntegration(t *testing.T) {
	token := getCircleAPIKey()

	client, err := circleci.NewClient(token, serverURL)
	require.Nil(t, err, "circleci.NewClient() error should be nil, received: %v", err)

	projects, err := client.Projects.List(context.Background())
	require.Nil(t, err, "client.Projects.List() error should be nil, received: %+v", errors.Unwrap(err))
	assert.GreaterOrEqual(t, len(projects), 1, "Expected at least 1 project in 'projects', got: %v", len(projects))
}

func TestGetProjectIntegration(t *testing.T) {
	token := getCircleAPIKey()

	client, err := circleci.NewClient(token, serverURL)
	require.Nil(t, err, "circleci.NewClient() error should be nil, received: %v", err)

	project, err := client.Projects.Get(context.Background(), proj, user)
	require.Nil(t, err, "client.Projects.Get() error should be nil, received: %+v", errors.Unwrap(err))
	assert.Equal(t, proj, project.Name, "project.Name returned unexpected value. Expected: %v. Got: %v", proj, project.Name)
}
