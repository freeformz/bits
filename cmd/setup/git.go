package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

func gitIgnoreLinesToAppend(fn string) ([]string, error) {
	lines := []string{
		"coverage.out",
		"coverage.html",
	}
	d, err := ioutil.ReadFile(fn)
	if err != nil {
		return lines, err
	}
	var add []string
	for _, l1 := range strings.Split(string(d), "\n") {
		var found bool
		for _, l2 := range lines {
			if l1 == l2 {
				found = true
			}
		}
		if !found {
			add = append(add, l1)
		}
	}
	// if there are lines to add and .gitignore doesn't end in a \n, add one first
	if len(add) > 0 && !bytes.HasSuffix(d, []byte("\n")) {
		add = append([]string{"\n"}, add...)
	}
	return add, nil
}

func setupGitIgnore(fn string) error {
	lines, err := gitIgnoreLinesToAppend(fn)
	if err != nil {
		if !os.IsNotExist(err) {
			return err
		}
	}
	f, err := os.OpenFile(fn, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	for i, l := range lines {
		if l != "\n" && i != len(lines)-1 {
			l += "\n"
		}
		if _, err = fmt.Fprint(f, l); err != nil {
			break
		}
	}
	return err
}
