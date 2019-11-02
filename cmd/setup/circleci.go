package main

import (
	"html/template"
	"os"
	"strings"
)

const (
	orbVersion       = "0.2.0"  //TODO[freeformz]: Figure this out instead of hard coding it
	golangCIVersion  = "1.20.0" //TODO[freeformz]: Figure this out instead of hard coding it
	circleCITemplate = `
version: 2.1
orbs:
  golang: heroku/golang@{{.OrbVersion}}

workflows:
  ci:
    jobs:
      - golang/golangci-lint:
          version: "v{{.GolangciLintVersion}}"
      - golang/test-nodb:
          version: "{{.GoVersion}}"
`
)

func createCircleCIConfig(cf, goVer string) error {
	t, err := template.New("circleci").Parse(circleCITemplate)
	if err != nil {
		return err
	}
	f, err := os.Create(cf)
	if err != nil {
		return err
	}
	defer f.Close()
	return t.Execute(f, struct {
		GoVersion, GolangciLintVersion, OrbVersion string
	}{strings.TrimPrefix(goVer, "go"), golangCIVersion, orbVersion})
}
