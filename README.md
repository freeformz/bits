# bits

## mage stuff

```console
$ mage -init
magefile.go created
```

magefile.go

```go
// +build mage

package main

import (
  // mage:import
  _ "github.com/freeformz/bits/mage/gons"
  // mage:import
  _ "github.com/freeformz/bits/mage/toolns"
)

```

```console
$ mage -f
Targets:
  go:cover                Run go tool cover, defaults: GO_COVER_ARGS="-html=coverage.out -o coverage.html"
  go:coverage             Open the coverage output in your browser (runs "go tool cover -html=coverage.out")
  go:test                 Run go test, defaults: GO_TEST_ARGS="-v -race -coverprofile=coverage.out -covermode=atomic ./..."
  golangcilint:install    golangci-lint to $TOOL_CACHE/bin/golangci-lint-<$GOLANGCILINT_VER>, defaults: GOLANGCILINT_VER=1.19.0.
  golangcilint:rm         Remove all cached versions of golangci-lint
  golangcilint:run        golangci-lint, defaults: GOLANGCILINT_VER=1.19.0.
```
