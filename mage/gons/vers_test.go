package gons

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/magefile/mage/sh"
)

func setBitCacheToTestData(t *testing.T) func() {
	t.Helper()
	f := func() {}
	orig, set := os.LookupEnv("BIT_CACHE")
	if set {
		f = func() {
			if err := os.Setenv("BIT_CACHE", orig); err != nil {
				t.Fatal("unexpected error")
			}
		}
	}
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal("unexpected error:", err)
	}
	td := filepath.Join(filepath.Dir(wd), "testdata")
	os.Setenv("BIT_CACHE", td)
	if err := filepath.Walk(td, func(path string, info os.FileInfo, err error) error {
		return sh.Run("touch", path)
	}); err != nil {
		t.Fatal("unexpected error:", err)
	}
	return f
}

func TestExpandVersion(t *testing.T) {
	cleanup := setBitCacheToTestData(t)
	defer cleanup()

	for name, tc := range map[string]struct {
		expected string
		err      error
	}{
		"go1.11":   {expected: "go1.11.13"}, //go1.11 selected because there is unlikely to be a new release of it.
		"go1.11.0": {expected: "go1.11"},
		"go1.11.x": {expected: "go1.11.13"},
	} {
		tc := tc
		name := name
		t.Run(name, func(t *testing.T) {
			got, err := expandVersion(name)
			if tc.err != err {
				t.Fatalf("expected %q, got %q", tc.err, err)
			}
			if tc.expected != got {
				t.Errorf("expected %q, got %q", tc.expected, got)
			}
		})
	}
}
