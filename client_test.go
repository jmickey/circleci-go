package circleci_test

import (
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"mickey.dev/go/circleci-go"
)

func getCircleAPIKey() string {
	return os.Getenv("CIRCLECI_KEY")
}

func TestNewClient(t *testing.T) {
	client, err := circleci.NewClient("TestAPIToken", "https://testserver.com")
	require.Nil(t, err, "Error should be nil")
	require.NotNil(t, client, "Client should no be nil")
	assert.Equal(t, "TestAPIToken", client.APIKey, "API Token doesn't match in client")
	assert.Equal(t, "https://testserver.com", client.ServerURL, "ServerURL doesn't match in client")
}

func TestNewClientWithCustomClient(t *testing.T) {
	c := &http.Client{Timeout: time.Minute * 5}
	client, err := circleci.NewClient("TestAPIToken", "https://testserver.com", circleci.SetBaseHTTPClient(c))
	require.Nil(t, err, "Error should be nil")
	require.NotNil(t, client, "Client should no be nil")
	assert.Equal(t, time.Minute*5, client.GetHTTPClient().Timeout, "Timeout should be 5 minutes")
}

func TestNewClientError(t *testing.T) {
	tests := []struct {
		token string
		url   string
	}{
		{"APIToken1", " Http://badURL1"},
		{"APIToken2", "%badURL2%"},
		{"APIToken3", "https:// badURL3"},
	}

	for _, test := range tests {
		_, err := circleci.NewClient(test.token, test.url)
		assert.NotNil(t, err, "Error for url '%v' should not be nil.", test.url)
	}
}
