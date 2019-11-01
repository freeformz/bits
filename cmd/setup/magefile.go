package main

import (
	"html/template"
	"os"
)

const mfTemplate = `// +build mage

package main

import (
    // mage:import
    _ "github.com/freeformz/bits/mage/toolns"
    // mage:import
    "github.com/freeformz/bits/mage/gons"
)

func init() {
    gons.ModuleName = "{{.Module}}"
    gons.Version = "{{.GoVersion}}"

    // Other gons settings (defaults)
    // gons.CoverArgs  = "-html=coverage.out -o coverage.html"
    // gons.TestArgs   = "-v -race -coverprofile=coverage.out -covermode=atomic ./..."

    // Golangci-lint settings (defaults), remove the '_' part of the import above
    // toolns.GolangciLint.Version = "1.20.0"
    // toolns.GolangciLint.RunArgs = "--fix ./..."
}`

func createMageFile(mfn, module, version string) error {
	t, err := template.New("magefile.go").Parse(mfTemplate)
	if err != nil {
		return err
	}
	f, err := os.Create(mfn)
	if err != nil {
		return err
	}
	defer f.Close()
	return t.Execute(f, struct {
		Module, GoVersion string
	}{module, version},
	)
}

func goModuleUpdate(f string) error {
	if err := run("go", "get", "github.com/freeformz/bits/mage"); err != nil {
		return err
	}
	if err := run("go", "mod", "tidy"); err != nil {
		return err
	}
	return run("go", "mod", "verify")
}
