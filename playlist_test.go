// Part of the Go playlist parser package
//
// Copyright 2013 Rafal Zajac rzajac<at>gmail<dot>com. All rights reserved.
// http://github.com/rzajac/plparser
//
// Licensed under the MIT license

package plparser

import (
	// "fmt"
	"testing"
)

func TestPlaylist(t *testing.T) {

	var files = map[string]struct {
		pltype    string
		detected  bool
		firstLine string
	}{
		"./testpls/asf1.asf":     {"asf", true, "[Reference]"},
		"./testpls/asf2.asf":     {"asf", true, "[Reference]"},
		"./testpls/asx1.asx":     {"asx", true, "<ASX version=\"3.0\" BANNERBAR=\"AUTO\">"},
		"./testpls/asx2.asx":     {"asx", true, "<ASX version=\"3.0\" BANNERBAR=\"AUTO\">"},
		"./testpls/m3u1.m3u":     {"m3u", true, "http://live1.example.com:2151/"},
		"./testpls/m3u2.m3u":     {"m3u", true, "http://live1.example.com:2151/"},
		"./testpls/pls1.pls":     {"pls", true, "[playlist]"},
		"./testpls/pls2.pls":     {"pls", true, "[playlist]"},
		"./testpls/pls3.pls":     {"pls", true, "[playlist]"},
		"./testpls/pls4.pls":     {"pls", true, "[playlist]"},
		"./testpls/unknown1.txt": {"", false, "Not a playlist"},
		"./testpls/unknown2.txt": {"", false, "[WeirdPlaylist]"},
	}

	for filePath, test := range files {
		f := getPLFile(filePath)

		plr := new(PlaylistResp)
		plr.Raw = f

		pl := NewPlaylist(plr)
		pltype, _ := pl.Parse()

		if test.pltype != pltype {
			t.Fatalf("Expected playlist %s to be of type '%s' but it's '%s'", filePath, test.pltype, pl.Type)
		}

		if test.pltype != pl.Type {
			t.Fatalf("Expected playlist %s to be of type '%s' but it's '%s'", filePath, test.pltype, pl.Type)
		}

		if test.firstLine != pl.firstLine {
			t.Fatalf("Expected first line to be '%s' and I got '%s' [%s]", test.firstLine, pl.firstLine, filePath)
		}

		if test.detected != pl.IsDetected() {
			t.Fatalf("Expected IsDetected() %v != %v [%s]", test.detected, pl.IsDetected(), filePath)
		}

	}

}
