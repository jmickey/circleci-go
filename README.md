# CircleCI Go Client
Go API Client for CircleCI. Currently implements the v1.1 API (https://circleci.com/docs/api).

**Currently under active development**

## Usage:

### Installation:

Change directory to you project dir and run:
```sh
go get mickey.dev/go/circleci-go@latest
```

Import the package:
```go
package somePackage

import (
    "log"
    "context"

    "mickey.dev/go/circleci-go"
)

const (
    API_TOKEN = "YourCircleCIAPIToken"
    SERVER_URL = "https://circleci.com"
)

func main() {
    client, err := circleci.NewClient(API_TOKEN, SERVER_URL)
    if err != nil {
        log.Fatalf("Couldn't create new CircleCI API Client: %v", err)
    }

    // Example call - List all followed projects:
    projects, _ := client.Projects.List(context.Background())
}
```

## Development

Currently this package is incomplete and under active development. PRs are welcome! 

To configure your build environment:

1. Fork the repo.
2. Clone: `git clone git@github.com/<YOUR_GITHUB_USERNAME>/circleci-go.git`
3. `cd circleci-go`
4. `go get`
5. `go test ./...` - Confirm the tests pass!

### Run Integration Tests

The integration tests will attempt to reach out to the live CircleCI API. As a result you'll need to retreive an API token from CircleCI and set the following environment variable:

```sh
export CIRCLECI_KEY=<YOUR_CIRCLECI_API_KEY>
```

You may wish to change the project the tests are run against, you can do this from within the [integration_test.go](/integration_test.go) file.

To run the tests use the following command:

```sh
go test -tags=integration
```

## To Do

- Everything.
- `/me` endpoints.
- Jobs and Workflows.