// +build integration

package circleci_test

import "os"

const (
	serverURL string = "https://circleci.com"
	proj      string = "circleci-go"
	ser       string = "jmickey"
)

func getCircleAPIKey() string {
	return os.Getenv("CIRCLECI_KEY")
}
