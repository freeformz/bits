package main

import (
	"io"
	"os/exec"
	"strings"
)

func run(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	return cmd.Run()
}

func output(name string, args ...string) (string, error) {
	cmd := exec.Command(name, args...)
	out, err := cmd.Output()
	return strings.TrimSpace(string(out)), err
}

func awk(s string, input io.Reader) (string, error) {
	cmd := exec.Command("awk", s)
	cmd.Stdin = input
	out, err := cmd.Output()
	return strings.TrimSpace(string(out)), err
}

func gofmt(t string) error {
	return run("gofmt", "-w", t)
}
