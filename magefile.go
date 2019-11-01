// +build mage

package main

import (
	// mage:import
	_ "github.com/freeformz/bits/mage/toolns"
	// mage:import
	"github.com/freeformz/bits/mage/gons"
)

func init() {
	gons.ModuleName = "github.com/freeformz/bits"
	gons.Version = "go1.13"

	// Other gons settings (defaults)
	// gons.CoverArgs  = "-html=coverage.out -o coverage.html"
	// gons.TestArgs   = "-v -race -coverprofile=coverage.out -covermode=atomic ./..."

	// Golangci-lint settings (defaults), remove the '_' part of the import above
	// toolns.GolangciLint.Version = "1.20.0"
	// toolns.GolangciLint.RunArgs = "--fix ./..."
}
