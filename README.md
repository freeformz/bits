# bits

### Quick Start

```console
$ export BITS_VERSION="$(curl -s https://api.github.com/repos/freeformz/bits/releases/latest | grep tag_name | awk '{ sub(/"/, "", $2); sub(/",/, "", $2); print $2 }')"
$ go mod init <module name>
$ go get github.com/freeformz/bits/cmd/setup@${BITS_VERSION}
$ git init .
$ git add go.mod go.sum
$ git commit -avm "new go module w/bits"
$ go run github.com/freeformz/bits/cmd/setup
go run github.com/freeformz/bits/cmd/setup

Repo directory: "/Users/emuller/devel/<project>"
Magefile: "/Users/emuller/devel/<project>/magefile.go"
go.mod: "/Users/emuller/devel/<project>/go.mod"
.golangci.tml: "/Users/emuller/devel/<project>/.golangci.yml"
.gitignore: "/Users/emuller/devel/<project>/.gitignore"

Module Name: "github.com/heroku/<project>"
Current major Go version: go1.13

setup complete!

Files to take a look at:
/Users/emuller/devel/<project>/magefile.go: Magefile used to configure bits / write additional tasks
/Users/emuller/devel/<project>/.golangci.yml: Golangci-lint config file
/Users/emuller/devel/<project>/.gitignore: git ignore file
/Users/emuller/devel/<project>/.circleci/config.yml: CircleCI config file

Don't forget to commit the changes/new files.

Run 'mage -f' to see the initial set of mage targets.
```

Creates a defaults magefile.go, .circleci config, .gitignore, & .golangci.yml

```console
$ mage -f
Targets:
  go:checkVersion        checks that the version of go being used is the version specified or the latest version
  go:cover               runs go tool cover with default args set from `CoverArgs`
  go:coverage            opens the coverage output in your browser (runs "go tool cover -html=coverage.out")
  go:test                runs `go test` with default args set from `TestArgs`
  golangcilint:remove    removes all cached versions of golangci-lint
  golangcilint:run       runs golangci-lint using RunArgs
```

### Targets

#### Go Namespace

`go:checkVersion` - Asserts that the version is use is the version specified. If the version specified ends in `.x` or only has two parts (`go1.13`) it is expanded to the most recent patch version of that go release. Modify the version by specifying `gons.Version` in your magefile.

`go:cover` - Generate cover file. Modify the arguments by specifying `gons.CoverArgs` in your magefile.

`go:coverage` - Generates coverage information and opens it in your browser. Modify the arguments by specifying `gons.CoverArgs` in your magefile.

`go:test` - Runs go test. Modify the arguments by specifying `gons.TestArgs` in your magefile.

#### Tool Namespace

`golangclilint:remove` - Removes all cached versions of golangci-lint.

`golangcilint:run` - Runs golangci-lint. Modify the arguments by specifying `toolns.GolangciLint.RunArgs` and change the version by specifying `toolns.GolangciLint.Version` in your magefile.