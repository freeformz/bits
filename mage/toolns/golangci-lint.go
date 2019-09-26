package tools

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

type Golangcilint mg.Namespace

const (
	// version,version,os,arch
	golangciURLFormat = "https://github.com/golangci/golangci-lint/releases/download/v%s/golangci-lint-%s-%s-%s.tar.gz"
	// version,os,arch
	golangciArchiveDirFormat = "golangci-lint-%s-%s-%s"
	golangciDefaultVersion   = "1.19.0"
)

// Install golangci-lint to $TOOL_CACHE/bin/golangci-lint-<$GOLANGCILINT_VER>, defaults: GOLANGCILINT_VER=1.19.0.
func (g Golangcilint) Install() error {
	t, err := g.path()
	if err != nil {
		return err
	}

	if _, err := os.Stat(t); os.IsNotExist(err) {
		ver := g.ver()
		if mg.Verbose() {
			fmt.Println("Downloading: golangci-lint-v" + ver + " to " + t)
		}
		if err := dlAndExtract(
			fmt.Sprintf(golangciURLFormat, ver, ver, runtime.GOOS, runtime.GOARCH),
			filepath.Join(fmt.Sprintf(golangciArchiveDirFormat, ver, runtime.GOOS, runtime.GOARCH), "golangci-lint"),
			t); err != nil {
			return err
		}
	}
	return nil
}

func (g Golangcilint) ver() string {
	ver := os.Getenv("GOLANGCILINT_VER")
	if ver == "" {
		ver = golangciDefaultVersion
	}
	return ver
}

func (g Golangcilint) path() (string, error) {
	d, err := binDir()
	if err != nil {
		return "", err
	}
	ver := os.Getenv("GOLANGCILINT_VER")
	if ver == "" {
		ver = golangciDefaultVersion
	}
	return filepath.Join(d, fmt.Sprintf("golangci-lint-%s", ver)), nil
}

func (g Golangcilint) runOpts() []string {
	opts := os.Getenv("GOLANGCILINT_RUN_OPTS")
	if opts == "" {
		return []string{"--fix", "./..."}
	}
	return strings.Split(opts, " ")
}

// Run golangci-lint, defaults: GOLANGCILINT_VER=1.19.0. GOLANGCILINT_RUN_OPTS="--fix ./..."
func (g Golangcilint) Run() error {
	mg.Deps(g.Install)
	p, err := g.path()
	if err != nil {
		return err
	}
	opts := append([]string{"run"}, g.runOpts()...)
	return sh.RunV(p, opts...)
}

// Remove all cached versions of golangci-lint
func (Golangcilint) Rm() error {
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
