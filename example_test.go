package circleci_test

import (
	"log"
	"net/http"

	"mickey.dev/go/circleci-go"
)

func ExampleWithBaseHTTPClient() {
	httpClient := &http.Client{}

	_, err := circleci.NewClient("YourCircleCIAPIToken", circleci.WithBaseHTTPClient(httpClient))
	if err != nil {
		log.Fatal(err)
	}
}
