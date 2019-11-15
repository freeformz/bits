package gons

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
	"github.com/magefile/mage/target"
)

var (
	// ModuleName if not set this is determined from the go.mod file
	ModuleName = moduleName()
	// Version of go to use. If not set it defaults to the latest version of Go
	Version = ""
	// CoverArgs to supply to go test.
	CoverArgs = "-html=coverage.out -o coverage.html"
	// TestArgs to supply to go test.
	TestArgs = "-v -race -coverprofile=coverage.out -covermode=atomic ./..."
)

//TODO: warning or error instead of just empty return
func moduleName() string {
	d, err := os.Getwd()
	if err != nil {
		return ""
	}
	f, err := os.Open(filepath.Join(d, "go.mod"))
	if err != nil {
		return ""
	}
	defer f.Close()
	s := bufio.NewScanner(f)
	var m string
	for s.Scan() {
		m = s.Text()
		p := strings.SplitN(m, " ", 2)
		if len(p) == 2 && p[0] == "module" {
			return p[1]
		}
	}
	return ""
}

// Go namespace
type Go mg.Namespace

var (
	goTest  = sh.OutCmd("go", "test")
	goCover = sh.RunCmd("go", "tool", "cover")
	goList  = sh.OutCmd("go", "list", "-json", "-find", "./...")
)

// Test runs `go test` with default args set from `TestArgs`
func (g Go) Test(ctx context.Context) error {
	mg.CtxDeps(ctx, g.CheckVersion)
	out, err := goTest(strings.Split(TestArgs, " ")...)
	if err != nil {
		fmt.Println(out)
	}
	return err
}

func goFiles() ([]string, error) {
	out, err := goList()
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

// CheckVersion checks that the version of go being used is the version specified or the latest version
func (g Go) CheckVersion(ctx context.Context) error {
	ver := Version
	if ver == "" {
		var err error
		ver, err = latestVersion()
		if err != nil {
			return err
		}
	}
	cv, err := sh.Output("go", "version")
	if err != nil {
		return err
	}
	scv := strings.Split(cv, " ")
	if len(scv) != 4 {
		return fmt.Errorf("unknown `go version` string: %q", cv)
	}
	ver, err = expandVersion(ver)
	if err != nil {
		return err
	}
	if ver != scv[2] {
		return fmt.Errorf("current version (%s) is not the same as specified/latest version (%s)", scv[2], ver)
	}
	fmt.Printf("current go version (%s) matches specified/latest version (%s)\n", scv[2], ver)
	return nil
}

// Cover runs go tool cover with default args set from `CoverArgs`
func (g Go) Cover(ctx context.Context) error {
	mg.CtxDeps(ctx, g.CheckVersion)
	gf, err := goFiles()
	if err != nil {
		return err
	}
	if need, _ := target.Path("coverage.out", gf...); need {
		mg.Deps(g.Test)
	}
	return goCover(strings.Split(CoverArgs, " ")...)
}

// Coverage opens the coverage output in your browser (runs "go tool cover -html=coverage.out")
func (g Go) Coverage(ctx context.Context) error {
	mg.CtxDeps(ctx, g.CheckVersion, g.Cover)
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
