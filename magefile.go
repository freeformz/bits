// +build mage

package main

import (
	// mage:import
	"github.com/freeformz/bits/mage/gons"
	// mage:import
	_ "github.com/freeformz/bits/mage/toolns"
)

func init() {
	gons.Version = "go1.13"
}
