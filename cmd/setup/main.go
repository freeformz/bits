package main

import (
	"fmt"
	"os"
	"path/filepath"
)

const dirMode = os.ModeDir | os.ModePerm

func main() { //nolint:gocyclo
	assertMageExecutes()

	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	cwd = assertGitTop(cwd)
	mf := filepath.Join(cwd, "magefile.go")
	gomodf := filepath.Join(cwd, "go.mod")
	golangcilintf := filepath.Join(cwd, ".golangci.yml")
	gif := filepath.Join(cwd, ".gitignore")
	ccid := filepath.Join(cwd, ".circleci")
	ccicf := filepath.Join(ccid, "config.yml")

	fmt.Println()
	fmt.Printf("Repo directory: %q\n", cwd)
	fmt.Printf("Magefile: %q\n", mf)
	fmt.Printf("go.mod: %q\n", gomodf)
	fmt.Printf(".golangci.tml: %q\n", golangcilintf)
	fmt.Printf(".gitignore: %q\n", gif)
	fmt.Println()

	assertGitPorcelain()
	assertFileExists(filepath.Join(cwd, "go.mod"), "No go.mod file. Use with go modules is required.")

	err = backupIfFileExists(mf)
	if err != nil && !os.IsNotExist(err) {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	if err == nil {
		fmt.Printf("Renamed existing %s to %s.bak\n", mf, mf)
	}

	modName, err := determineModuleName(gomodf)
	if err != nil || modName == "" {
		fmt.Printf("unable to determine go module name: %s\n", err.Error())
		os.Exit(1)
	}
	fmt.Printf("Module Name: %q\n", modName)

	goVersion, err := determineCurrentGoVersion()
	if err != nil {
		fmt.Printf("unable to determine current go version: %s\n", err.Error())
		os.Exit(1)
	}
	goVersion = truncateGoVersion(goVersion)
	fmt.Printf("Current major Go version: %s\n", goVersion)

	// magefile.go
	if err := createMageFile(mf, modName, goVersion); err != nil {
		fmt.Printf("unable to create %s: %s\n", mf, err.Error())
		os.Exit(1)
	}

	if err := gofmt(mf); err != nil {
		fmt.Printf("unable to gofmt %s: %s\n", mf, err.Error())
		os.Exit(1)
	}

	if err := goModuleUpdate(); err != nil {
		fmt.Printf("unable to prep go.mod file: %s\n", err.Error())
		os.Exit(1)
	}

	// golangci-lint
	err = backupIfFileExists(golangcilintf)
	if err != nil && !os.IsNotExist(err) {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	if err == nil {
		fmt.Printf("Renamed existing %s to %s.bak\n", golangcilintf, golangcilintf)
	}
	if err := createGolangcilintConfig(golangcilintf, modName); err != nil {
		fmt.Printf("unable to create %q: %s\n", golangcilintf, err.Error())
		os.Exit(1)
	}

	//.gitignore
	if err := setupGitIgnore(gif); err != nil {
		fmt.Printf("unable to setup %q: %s\n", gif, err.Error())
		os.Exit(1)
	}

	// CircleCI
	if err := os.MkdirAll(ccid, dirMode); err != nil {
		fmt.Printf("unable to create circle ci directory %q: %s\n", ccid, err.Error())
		os.Exit(1)
	}
	err = backupIfFileExists(ccicf)
	if err != nil && !os.IsNotExist(err) {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	if err == nil {
		fmt.Printf("Renamed existing %s to %s.bak\n", ccicf, ccicf)
	}
	if err := createCircleCIConfig(ccicf, goVersion); err != nil {
		fmt.Printf("unable to setup circle ci config %q: %s\n", ccicf, err.Error())
		os.Exit(1)
	}

	fmt.Printf(`
setup complete!

Files to take a look at:
%s: Magefile used to configure bits / write additional tasks
%s: Golangci-lint config file
%s: git ignore file
%s: CircleCI config file

Don't forget to commit the changes/new files.

Run 'mage -f' to see the initial set of mage targets

`, mf, golangcilintf, gif, ccicf)
}

func determineModuleName(mf string) (string, error) {
	i, err := os.Open(mf)
	if err != nil {
		return "", err
	}
	defer i.Close()
	return awk(`/^module/ { print $2; exit}`, i)
}

func backupIfFileExists(f string) error {
	fi, err := os.Stat(f)
	if os.IsNotExist(err) {
		return err
	}
	if err != nil {
		return fmt.Errorf("unexpected error backing up file (%s) if it exists: %s", f, err.Error())
	}
	if fi.IsDir() {
		return fmt.Errorf("%q is a directory, not a file", f)
	}
	return os.Rename(f, f+".bak")
}

// file name, message if file doesn't exist.
func assertFileExists(f, msg string) {
	_, err := os.Stat(f)
	if os.IsNotExist(err) {
		fmt.Println(msg)
		os.Exit(1)
	}
	if err != nil {
		fmt.Printf("Unexpected error asserting file (%s) exists:\n%s\n", f, err.Error())
		os.Exit(1)
	}
}

func assertGitPorcelain() {
	out, err := output("git", "status", "--porcelain")
	if err != nil {
		fmt.Printf("Unable to execute 'git status --porcelain':\n%s\n", err.Error())
		os.Exit(1)
	}
	if len(out) != 0 {
		fmt.Println(`'git status' shows there are change that need to be committed.
'setup' should be run on a clean git repo so that it's easy to evaluate and/or throw away changes.
Please commit, stash, or throw away pending changes, then re-run setup.`)
		os.Exit(1)
	}
}

func assertGitTop(top string) string {
	out, err := output("git", "rev-parse", "--show-toplevel")
	if err != nil {
		fmt.Println("Unable to execute 'git rev-parse --show-toplevel':")
		fmt.Println(err)
		os.Exit(1)
	}
	if out != top {
		fmt.Printf("Current directory (%s) is not the root of the git repo (%s)\nChanging directory to %s\n", top, out, out)
		if err := os.Chdir(out); err != nil {
			fmt.Printf("Unable to change directory to %q:\n%s\n", out, err.Error())
			os.Exit(1)
		}
	}
	return out
}

func assertMageExecutes() {
	_, err := output("mage", "-version")
	if err != nil {
		fmt.Println("Mage is required")
		fmt.Println("`brew install mage` (on macOS) or follow the instructions here: https://magefile.org/")
		os.Exit(1)
	}
}
