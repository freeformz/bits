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
    echo "`brew install mage` (on macOS) or follow instructions here: https://magefile.org/"
    exit 1
fi

if [ "${CWD}" != "${GIT_TOP}" ]; then
    echo "Current directory (${CWD}) is not the root of this git repo (${GIT_TOP})"
    echo "Please change to ${GIT_TOP} and re-run this script"
    exit 1
fi

if [ ! -z "$(git status --porcelain)" ]; then
    echo "'git status' shows there are changes that need to be committed."
    echo "'setup.sh' should be run on a clean git repo so that it's easy to evaluate and/or throw away changes."
    echo "Please commit, stash, or throw away pending changes."
    echo
    git status
    exit 1
fi

if [ ! -e "${GOMOD}" ]; then
    echo "Current directory does not contain a go.mod file"
    echo "use with go modules is required"
    exit 1
fi

#
# mage
#

if [ -e "${MF}" ]; then
    echo "Current directory contains a magefile.go file already, renaming to magefile.go.bak"
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
    gons.ModuleName = "${MODULE}"
    gons.Version = "${CGV}"

    // Other gons settings (defaults)
    // gons.CoverArgs  = "-html=coverage.out -o coverage.html"
    // gons.TestArgs   = "-v -race -coverprofile=coverage.out -covermode=atomic ./..."

    // Golangci-lint settings (defaults), remove the '_' part of the import above
    // toolns.GolangciLint.Version = "1.20.0"
    // toolns.GolangciLint.RunArgs = "--fix ./..."
}
EOF

gofmt -w ${MF}
go get github.com/freeformz/bits/mage@v0.0.3
#go mod edit -replace=github.com/freeformz/bits=../bits

#
# golangci-lint
# TODO: Move to a mage task
#
writeGolangcilintConfig() {
    local f=${1}
    echo
    echo "creating ${f}"
    cat << EOF > ${f}
# See https://github.com/golangci/golangci-lint#config-file
run:
  deadline: 1m #Default
  issues-exit-code: 1 #Default
  tests: true #Default

linters:
  enable:
    - deadcode
    - dupl
    - errcheck
    - gochecknoinits
    - goconst
    - gocritic
    - gocyclo
    - goimports
    - golint
    - gosec
    - gosimple
    - govet
    - ineffassign
    - lll
    - maligned
    - misspell
    - nakedret
    - prealloc
    - scopelint
    - staticcheck
    - structcheck
    - stylecheck
    - typecheck
    - unconvert
    - unparam
    - unused
    - varcheck

issues:
  exclude-rules:
    - path: _test\.go
      linters:
        - goconst # Don't run on tests because they often repeate the same string
        - lll # Don't do line length checks in test code.
        - dupl # Sometimes tests duplicate for the sake of clarity.

linters-settings:
  misspell:
    locale: US
    #ignore-words:
    #  - someword
  goimports:
    local-prefixes: ${MODULE}
  gocyclo:
    # minimal code complexity to report, 30 by default (but we recommend 10-20)
    min-complexity: 15
  lll:
    # max line length, lines longer will be reported. Default is 120.
    line-length: 130
  maligned:
    suggest-new: true
EOF
}

GLCLF="${CWD}/.golangci.yml"
if [ -e "${GLCLF}" ]; then
    echo
    echo "golangci-lint config file exists (${GLCLF})"
    echo "Do you want to back it up and replace it with the bits default?"
    while true; do
        read -p 'Y (backup) / N (overwrite): ' response
        if [ "${response}" == "Y" -o "${response}" == "N" -o "${response}" == "y" -o "${response}" == "n" ]; then
            break
        fi
    done
    if [ "${response}" == "Y" -o "${response}" == "y" ]; then
        echo
        echo "Moving .golangci.yml to .golanci.yml.bak"
        mv ${GLCLF} ${GLCLF}.bak
        writeGolangcilintConfig ${GLCLF}
      fi
else
    writeGolangcilintConfig ${GLCLF}
fi

#
# .gitignore
# TODO: Move to a mage task
#
GI="${CWD}/.gitignore"
if [ ! -e "${GI}" ]; then
    touch ${GI}
fi

if [ -z "$(grep "^coverage.out$" < ${GI})" ]; then
    echo "coverage.out" >> ${GI}
fi

if [ -z "$(grep "^coverage.html$" < ${GI})" ]; then
    echo "coverage.html" >> ${GI}
fi

#
# CircleCI
# TODO: Move to a mage task
#
CID="${CWD}/.circleci"
CIC="${CID}/config.yml"
mkdir -p {$CID}

writeCircleCIConfig() {
    local f=${1}
    echo
    echo "creating ${f}"
    cat << EOF > ${f}
version: 2.1
orbs:
  golang: heroku/golang@0.2.0

workflows:
  ci:
    jobs:
      - golang/golangci-lint:
          version: "v1.20.0"
      - golang/test-nodb:
          version: "1.13"
EOF
}

if [ -e "${CIC}" ]; then
    echo
    echo "circle ci config file exists (${CIC})"
    echo "Do you want to back it up and replace it with the bits default?"
    while true; do
        read -p 'Y (backup) / N (overwrite): ' response
        if [ "${response}" == "Y" -o "${response}" == "N" -o "${response}" == "y" -o "${response}" == "n" ]; then
              break
        fi
      done
      if [ "${response}" == "Y" -o "${response}" == "y" ]; then
          echo
          echo "Moving .golangci.yml to .golanci.yml.bak"
          mv ${CIC} ${CIC}.bak
        writeCircleCIConfig ${CIC}
      fi
else
    writeCircleCIConfig ${CIC}
fi

echo
echo "setup complete"
echo
echo "Files to take a look at:"
echo -e "\t${MF}"
echo -e "\t${GLCLF}"
echo -e "\t${GI}"
echo -e "\t${CIC}"
echo
echo "Don't forget to commit the changes/new files."
echo
mage -f