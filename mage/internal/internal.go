package internal

import (
	"fmt"
	"io"
	"os"
	"os/user"
	"path/filepath"
	"time"
)

const (
	maxAge = 24 * time.Hour
)

var (
	dir = ".bit"

	ErrMaxAge = fmt.Errorf("file is older than 24 hours")
)

const dirMode = os.ModeDir | os.ModePerm

// DefaultCacheDirectory to use. Default is $HOME/.bit
func DefaultCacheDirectory() (string, error) {
	user, err := user.Current()
	if err != nil {
		return "", err
	}
	return filepath.Join(user.HomeDir, dir), nil
}

// CacheDirectory to use for the given namespace (ns). Uses $BIT_CACHE as the base if set.
func CacheDirectory(ns string) (string, error) {
	tl := os.Getenv("BIT_CACHE")
	var err error
	if tl == "" {
		tl, err = DefaultCacheDirectory()
	}
	if err != nil {
		return tl, err
	}
	if ns != "" {
		tl = filepath.Join(tl, ns)
	}
	return tl, os.MkdirAll(tl, dirMode)
}

// CachedFile to use for the given namespace (ns) and file name (ns). Uses CacheDirectory() as the base.
// May return ErrMaxAge as an indicator to the caller that the file is older than maxAge
func CachedFile(ns, fn string) (io.ReadCloser, error) {
	// locate the file based on the cache dir
	cd, err := CacheDirectory(ns)
	if err != nil {
		return nil, err
	}
	fp := filepath.Join(cd, fn)

	// if it's not there or older than maxAge return and error
	fi, err := os.Stat(fp)
	if err != nil {
		return nil, err
	}
	if time.Since(fi.ModTime()) > maxAge {
		return nil, ErrMaxAge
	}
	return os.Open(fp)
}
