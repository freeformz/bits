#!/usr/bin/env bash

set -e

#
# curl https://raw.githubusercontent.com/freeformz/bits/master/mage/setup.sh | bash
#

CWD=$(pwd)
GIT_TOP=$(git rev-parse --show-toplevel)
MF="${CWD}/magefile.go"
GOMOD="${CWD}/go.mod"

if [ -z "$(which mage)" ]; then
  echo "Mage is required"
	echo "`brew install mage` or follow instructions here: https://magefile.org/"
	exit 1
fi

if [ "${CWD}" != "${GIT_TOP}" ]; then
  echo "Current directory (${CWD}) is not the root of this git repo (${GIT_TOP})"
	echo "Please change to ${GIT_TOP} and re-run this script"
	exit 1
fi

if [ ! -e "${GOMOD}" ]; then
  echo "Current directory does not contain a go.mod file"
	echo "use with go modules is required"
	exit 1
fi

if [ -e "${MF}" ]; then
  echo "Current directory contains a magefile.go already"
	echo "renaming to magefile.go.bak"
	mv ${MF} ${MF}.bak
fi

MODULE=$(awk '/^module/ { print $2; exit }' < go.mod)
if [ -z "${MODULE}" ]; then
  echo "Unable to determine go module name"
	exit 1
fi

CGV=$(curl -s "https://golang.org/VERSION?m=text" | awk 'BEGIN { FS="."} { printf "%s.%s",$1,$2; exit }')

if [ -z "${CGV}" ]; then
  echo "Unable to determine current go version"
	exit 1
fi

cat << EOF > ${MF}
// +build mage

package main

import (
	// mage:import
	_ "github.com/freeformz/bits/mage/toolns"
	// mage:import
	"github.com/freeformz/bits/mage/gons"
)

func init() {
	// Defaults to whatever go.mod says, but can be overridden here.
	gons.ModuleName = "${MODULE}"
	// Defaults to https://golang.org/VERSION?m=text output, but can be overridden here
	gons.Version = "${CGV}"
}
EOF

gofmt -w ${MF}

go get github.com/freeformz/bits/mage
go mod edit -replace=github.com/freeformz/bits=../bits
echo "setup complete"
echo "Take a look at the values set in magefile.go"
mage -f