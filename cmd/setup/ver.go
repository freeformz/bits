package main

import (
	"io/ioutil"
	"net/http"
	"strings"
)

func truncateGoVersion(ver string) string {
	ver = strings.TrimPrefix(ver, "go")
	p := strings.Split(ver, "rc")
	ver = p[0]
	p = strings.Split(ver, "beta")
	ver = p[0]
	ver = strings.TrimSuffix(ver, "rc")
	p = strings.Split(ver, ".")
	return "go" + p[0] + "." + p[1]
}

func determineCurrentGoVersion() (string, error) {
	resp, err := http.Get("https://golang.org/VERSION?m=text")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	d, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(d), nil
}
