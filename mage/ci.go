package mage

import (
	"github.com/magefile/mage/mg"
)

type CI mg.Namespace

// hello
func (CI) LocalFoo() error {
	return nil
}
