package tools

import (
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/user"
	"path/filepath"

	"github.com/magefile/mage/sh"
	"github.com/mholt/archiver"
)

const dirMode = os.ModeDir | os.ModePerm

var defaultCache = filepath.Join(".bit", "toolcache")

func defaultToolDir() (string, error) {
	user, err := user.Current()
	if err != nil {
		return "", err
	}
	return filepath.Join(user.HomeDir, defaultCache), nil
}

func toolsDir() (string, error) {
	tl := os.Getenv("TOOL_CACHE")
	if tl == "" {
		return defaultToolDir()
	}
	return tl, nil
}

func binDir() (string, error) {
	d, err := toolsDir()
	if err != nil {
		return "", err
	}
	d = filepath.Join(d, "bin")
	err = os.MkdirAll(d, dirMode)
	if os.IsExist(err) {
		err = nil
	}
	return d, err
}

func repoRoot() (string, error) {
	return sh.Output("git", "rev-parse", "--show-toplevel")
}

// TODO: Support something other than .tar.gz
func dlAndExtract(source, file, dest string) error {
	td, err := ioutil.TempDir("", "x-mage")
	if err != nil {
		return err
	}
	defer func() { _ = os.RemoveAll(td) }()

	f, err := os.Create(filepath.Join(td, filepath.Base(source)))
	if err != nil {
		return err
	}
	defer func() { _ = f.Close() }()

	r, err := http.Get(source) //nolint:gosec
	if err != nil {
		return err
	}
	defer r.Body.Close()

	if _, err := io.Copy(f, r.Body); err != nil {
		return err
	}
	f.Close()

	if err := archiver.Extract(f.Name(), file, td); err != nil {
		return err
	}

	df, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer df.Close()

	sf, err := os.Open(filepath.Join(td, file))
	if err != nil {
		return err
	}
	defer sf.Close()

	if _, err = io.Copy(df, sf); err != nil {
		return err
	}
	return os.Chmod(df.Name(), 0755)
}
