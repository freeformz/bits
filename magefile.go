// +build mage

package main

import (
	// mage:import
	"github.com/freeformz/bits/mage/gons"
)

func init() {
	gons.Version = "go1.13"
}
