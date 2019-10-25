package toolns

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

// Golangcilint mage namespace
type Golangcilint mg.Namespace

const (
	// version,version,os,arch
	golangciURLFormat = "https://github.com/golangci/golangci-lint/releases/download/v%s/golangci-lint-%s-%s-%s.tar.gz"
	// version,os,arch
	golangciArchiveDirFormat = "golangci-lint-%s-%s-%s"
)

var (
	// Golangci configuration
	Golangci = struct {
		Version string
		RunArgs string
	}{
		Version: "1.20.0",
		RunArgs: "--fix ./...",
	}
)

// Install installs golangci-lint to $BIT_CACHE/tools/bin/golangci-lint-<Version>
func (g Golangcilint) Install() error {
	t, err := g.path()
	if err != nil {
		return err
	}

	if _, err := os.Stat(t); os.IsNotExist(err) {
		if mg.Verbose() {
			fmt.Println("Downloading: golangci-lint-v" + Golangci.Version + " to " + t)
		}
		if err := dlAndExtract(
			fmt.Sprintf(golangciURLFormat, Golangci.Version, Golangci.Version, runtime.GOOS, runtime.GOARCH),
			filepath.Join(fmt.Sprintf(golangciArchiveDirFormat, Golangci.Version, runtime.GOOS, runtime.GOARCH), "golangci-lint"),
			t); err != nil {
			return err
		}
	}
	return nil
}

func (g Golangcilint) path() (string, error) {
	d, err := binDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(d, fmt.Sprintf("golangci-lint-%s", Golangci.Version)), nil
}

// Run runs golangci-lint using RunArgs
func (g Golangcilint) Run() error {
	mg.Deps(g.Install)
	p, err := g.path()
	if err != nil {
		return err
	}
	opts := append([]string{"run"}, strings.Split(Golangci.RunArgs, " ")...)
	return sh.RunV(p, opts...)
}

// Remove  removes all cached versions of golangci-lint
func (Golangcilint) Remove() error {
	d, err := binDir()
	if err != nil {
		return err
	}

	m, err := filepath.Glob(filepath.Join(d, "golangci-lint*"))
	if err != nil {
		return err
	}
	for _, f := range m {
		if err := os.Remove(f); err != nil {
			return err
		}
	}
	return nil
}
