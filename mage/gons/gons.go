package gons

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
	"github.com/magefile/mage/target"
)

var (
	ModuleName = moduleName() // ModuleName, if not set this is determined from the go.mod file
	Version    = version()    // Version of go to use, if not set this defaults to what is returned from `https://golang.org/VERSION?m=text`
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
			fmt.Println("ModuleName:", p[1])
			return p[1]
		}
	}
	return ""
}

//TODO: warning or error instead of just empty return
func version() string {
	r, err := http.Get("https://golang.org/VERSION?m=text")
	if err != nil {
		return ""
	}
	d, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return ""
	}
	fmt.Println("GoVersion:", string(d))
	return string(d)
}

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
