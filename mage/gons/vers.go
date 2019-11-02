package gons

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"

	"github.com/freeformz/bits/mage/internal"
)

var (
	ns                   = "gons"
	versionCacheFileName = "versions.json"
)

func fetchAndCacheVersions() ([]string, error) {
	cd, err := internal.CacheDirectory(ns)
	if err != nil {
		return nil, err
	}
	f, err := os.Create(filepath.Join(cd, versionCacheFileName))
	if err != nil {
		return nil, err
	}
	defer f.Close()
	suffix := "." + runtime.GOOS + "-" + runtime.GOARCH + ".tar.gz"
	var vers []string
	var marker string
	for {
		url := "https://storage.googleapis.com/golang/?prefix=go"
		if marker != "" {
			url += "&marker=" + marker
		}
		r, err := http.Get(url)
		if err != nil {
			return nil, err
		}
		d := xml.NewDecoder(r.Body)
		rd := struct {
			IsTruncated bool
			NextMarker  string
			Contents    []struct {
				Key  string
				ETag string
			}
		}{}
		if err := d.Decode(&rd); err != nil {
			return nil, err
		}
		for _, v := range rd.Contents {
			if strings.HasSuffix(v.Key, suffix) {
				vers = append(vers, strings.TrimSuffix(v.Key, suffix))
			}
		}
		if !rd.IsTruncated {
			break
		}
		marker = rd.NextMarker
	}
	e := json.NewEncoder(f)
	return vers, e.Encode(vers)
}

func versions() ([]string, error) {
	cd, err := internal.CachedFile(ns, versionCacheFileName)
	if err != nil {
		if err == internal.ErrMaxAge || os.IsNotExist(err) {
			return fetchAndCacheVersions()
		}
		return nil, err
	}
	defer cd.Close()
	d := json.NewDecoder(cd)
	var v []string
	return v, d.Decode(&v)
}

//TODO[freeformz]: Refactor and remove nolint
//nolint:gocyclo
func expandVersion(ver string) (string, error) {
	p := strings.Split(ver, ".")
	if len(p) == 3 && p[2] != "x" { // already expanded, except if the patch is ".x"
		if p[2] == "0" { // if the patch is ".0", though it's the original version
			return fmt.Sprintf("%s.%s", p[0], p[1]), nil
		}
		return ver, nil
	}
	if len(p) == 3 && p[2] == "x" { // we want any version in the series, so truncate off so the reduction can happen below
		ver = fmt.Sprintf("%s.%s", p[0], p[1])
	}
	vers, err := versions()
	if err != nil {
		return "", err
	}
	// filter out impossible versions via a prefix match
	var match []string
	for _, v := range vers {
		if strings.HasPrefix(v, ver) {
			match = append(match, v)
		}
	}
	var majv, minv, pv int
	for _, v := range match {
		v = strings.TrimPrefix(v, "go")
		p := strings.Split(v, ".")
		t, err := strconv.Atoi(p[0])
		if err != nil {
			return "", err
		}
		if t > majv {
			majv = t
			minv = 0
			pv = 0
		}
		if strings.Contains(p[1], "beta") || strings.Contains(p[1], "rc") {
			// TODO: Fix beta / rc handling
		} else {
			t, err = strconv.Atoi(p[1])
			if err != nil {
				return "", err
			}
			if t > minv {
				minv = t
				pv = 0
			}
		}
		if len(p) >= 3 {
			t, err := strconv.Atoi(p[2])
			if err != nil {
				return "", err
			}
			if t > pv {
				pv = t
			}
		}
	}
	ver = fmt.Sprintf("go%d.%d", majv, minv)
	if pv > 0 {
		ver += fmt.Sprintf(".%d", pv)
	}
	return ver, nil
}

func latestVersion() (string, error) {
	r, err := http.Get("https://golang.org/VERSION?m=text")
	if err != nil {
		return "", err
	}
	d, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return "", err
	}
	return string(d), nil
}
