package tizzy

// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

import (
	"archive/zip"
	"errors"
	"io/ioutil"
	"log"
	"runtime"
	"time"
)

var (
	locations = map[string]time.Location{}
)

func init() {
	zoneinfo := runtime.GOROOT() + "/lib/time/zoneinfo.zip"

	r, err := zip.OpenReader(zoneinfo)
	if err != nil {
		log.Fatalf("Could not open zoneinfo.zip for reading: path '%s' error: '%v'", zoneinfo, err)
	}
	defer r.Close()

	for _, f := range r.File {
		if f == nil {
			log.Fatalf("zoneinfo.zip contained nil file pointer")
		}

		if f.FileHeader.FileInfo().IsDir() {
			continue
		}

		if _, ok := locations[f.Name]; ok {
			continue
		}

		bd, err := f.Open()
		if err != nil {
			log.Fatalf("Could not open file '%s': error: '%v'", f.Name, err)
		}

		tz, err := ioutil.ReadAll(bd)
		bd.Close()
		if err != nil {
			log.Fatalf("Could not read file '%s': error: '%v'", f.Name, err)
		}

		loc, err := time.LoadLocationFromTZData(f.Name, tz)
		if err != nil {
			log.Fatalf("Could not load location '%s': error: '%v'", f.Name, err)
		}

		locations[f.Name] = *loc
	}
}

// LoadLocation returns the Location with the given name.
//
// If the name is "" or "UTC", LoadLocation returns UTC.
// If the name is "Local", LoadLocation returns Local.
//
// Otherwise, the name is taken to be a location name corresponding to a file
// in the IANA Time Zone database, such as "America/New_York".
func LoadLocation(name string) (*time.Location, error) {
	l, err := LoadLocationValue(name)
	if err != nil {
		return nil, err
	}
	return &l, err
}

// LoadLocationValue avoids allocations on the heap.
func LoadLocationValue(name string) (time.Location, error) {
	if name == "" || name == "UTC" {
		return *time.UTC, nil
	}

	if name == "Local" {
		return *time.Local, nil
	}

	if containsDotDot(name) || name[0] == '/' || name[0] == '\\' {
		// No valid IANA Time Zone name contains a single dot,
		// much less dot dot. Likewise, none begin with a slash.
		return time.Location{}, errors.New("time: invalid location name")
	}

	zoneData, ok := locations[name]
	if !ok {
		return time.Location{}, errors.New("unknown time zone " + name)
	}

	return zoneData, nil
}

// containsDotDot reports whether s contains "..".
func containsDotDot(s string) bool {
	if len(s) < 2 {
		return false
	}
	for i := 0; i < len(s)-1; i++ {
		if s[i] == '.' && s[i+1] == '.' {
			return true
		}
	}
	return false
}
