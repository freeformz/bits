package gons

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
	"github.com/magefile/mage/target"
)

type Go mg.Namespace

var (
	goTest           = sh.RunCmd("go", "test")
	goCover          = sh.RunCmd("go", "tool", "cover")
	DefaultTestArgs  = []string{"-v", "-race", "-coverprofile=coverage.out", "-covermode=atomic", "./..."}
	DefaultCoverArgs = []string{"-html=coverage.out", "-o", "coverage.html"}
)

func testArgs() []string {
	e := os.Getenv("GO_TEST_ARGS")
	if e == "" {
		return DefaultTestArgs
	}
	return strings.Split(e, " ")
}

// Run go test, defaults: GO_TEST_ARGS="-v -race -coverprofile=coverage.out -covermode=atomic ./..."
func (g Go) Test() error {
	return goTest(testArgs()...)
}

func coverArgs() []string {
	e := os.Getenv("GO_COVER_ARGS")
	if e == "" {
		return DefaultCoverArgs
	}
	return strings.Split(e, " ")
}

func goFiles() ([]string, error) {
	out, err := sh.Output("go", "list", "-json", "-find", "./...")
	if err != nil {
		return nil, err
	}
	type glp struct {
		Dir     string
		GoFiles []string
	}
	b := strings.NewReader(out)
	d := json.NewDecoder(b)
	var goFiles []string
	for d.More() {
		var t glp
		if err := d.Decode(&t); err != nil {
			return goFiles, err
		}
		for _, f := range t.GoFiles {
			goFiles = append(goFiles, filepath.Join(t.Dir, f))
		}
	}
	return goFiles, nil
}

// Run go tool cover, defaults: GO_COVER_ARGS="-html=coverage.out -o coverage.html"
func (g Go) Cover() error {
	gf, err := goFiles()
	if err != nil {
		return err
	}
	need, _ := target.Path("coverage.out", gf...)
	if need {
		mg.Deps(g.Test)
	}
	return goCover(coverArgs()...)
}

// Open the coverage output in your browser (runs "go tool cover -html=coverage.out")
func (g Go) Coverage() error {
	gf, err := goFiles()
	if err != nil {
		return err
	}
	need, _ := target.Path("coverage.out", gf...)
	if need {
		mg.Deps(g.Test)
	}
	return goCover("-html=coverage.out")
}
