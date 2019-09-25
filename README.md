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
	// mage:import tool
	_ "github.com/freeformz/bits/mage/tools"
)
```

```console
$ mage
Targets:
  tool:golangcilint:install    golangci-lint to $TOOL_CACHE/bin/golangci-lint-<$GOLANGCILINT_VER), default GOLANGCILINT_VER=1.19.0.
  tool:golangcilint:rm         Remove all cached versions of golangci-lint
  tool:golangcilint:run        golangci-lint.
```
