package main

import (
	"html/template"
	"os"
)

const golangciTemplate = `# See https://github.com/golangci/golangci-lint#config-file
run:
  deadline: 1m #Default
  issues-exit-code: 1 #Default
  tests: true #Default

linters:
  enable:
    - deadcode
    - dupl
    - errcheck
    - gochecknoinits
    - goconst
    - gocritic
    - gocyclo
    - goimports
    - golint
    - gosec
    - gosimple
    - govet
    - ineffassign
    - lll
    - maligned
    - misspell
    - nakedret
    - prealloc
    - scopelint
    - staticcheck
    - structcheck
    - stylecheck
    - typecheck
    - unconvert
    - unparam
    - unused
    - varcheck

issues:
  exclude-rules:
    - path: _test\.go
      linters:
        - goconst # Don't run on tests because they often repeat the same string
        - lll # Don't do line length checks in test code.
        - dupl # Sometimes tests duplicate for the sake of clarity.

linters-settings:
  misspell:
    locale: US
    #ignore-words:
    #  - someword
  goimports:
    local-prefixes: {{.Module}}
  gocyclo:
    # minimal code complexity to report, 30 by default (but we recommend 10-20)
    min-complexity: 15
  lll:
    # max line length, lines longer will be reported. Default is 120.
    line-length: 130
  maligned:
    suggest-new: true`

func createGolangcilintConfig(configFile, module string) error {
	t, err := template.New("golangci.yml").Parse(golangciTemplate)
	if err != nil {
		return err
	}
	f, err := os.Create(configFile)
	if err != nil {
		return err
	}
	defer f.Close()
	return t.Execute(f, struct{ Module string }{module})
}
