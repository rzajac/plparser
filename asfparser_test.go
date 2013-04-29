// Part of the Go playlist parser package
//
// Copyright 2013 Rafal Zajac rzajac<at>gmail<dot>com. All rights reserved.
// http://github.com/rzajac/plparser
//
// Licensed under the MIT license

package plparser

import (
	// "regexp"
	"testing"
)

func TestAsfRegExp(t *testing.T) {

	var testUrlLines = []struct {
		line     string
		expected string
		idx      string
		matches  int
	}{
		{"Ref1=http://live.str3am.com:2150/", "http://live.str3am.com:2150/", "1", 3},
		{"ReF2=http://live.str3am.com:2150/", "http://live.str3am.com:2150/", "2", 3},
		{"REF1 = http://live.str3am.com:2150/", "http://live.str3am.com:2150/", "1", 3},
		{"Ref2 =http://live.str3am.com:2150/", "http://live.str3am.com:2150/", "2", 3},
		{"Ref2= http://live.str3am.com:2150/", "http://live.str3am.com:2150/", "2", 3},
		{"", "", "0", 0},
		{"not valid line", "", "0", 0},
	}

	for _, testLine := range testUrlLines {
		values := asfReg.FindStringSubmatch(testLine.line)

		if len(values) != testLine.matches {
			t.Fatalf("Expected to find %d matches for %s", testLine.matches, testLine.line)
		}

		if testLine.matches > 0 {

			if values[2] != testLine.expected {
				t.Fatalf("Expected url '%s' == '%s'", testLine.expected, values[2])
			}

			if values[1] != testLine.idx {
				t.Fatalf("Expected idx '%s' == '%s'", testLine.expected, values[1])
			}
		}
	}
}

func TestAsfFiles(t *testing.T) {

	var asf1 = map[int]plTestStruct{
		0: {},
		1: {1, "", "", "", "", "", "", "http://live.example.com/aaa?MSWMExt=.asf"},
		2: {2, "", "", "", "", "", "", "http://live2.example.com/bbb?MSWMExt=.asf"},
	}

	var asf2 = map[int]plTestStruct{
		0: {},
		1: {1, "", "", "", "", "", "", "http://live.example.com/aaa?MSWMExt=.asf"},
		2: {2, "", "", "", "", "", "", "http://live2.example.com/bbb?MSWMExt=.asf"},
		3: {3, "", "", "", "", "", "", "http://live3.example.com/ccc?MSWMExt=.asf"},
	}

	var testFiles = map[string]map[int]plTestStruct{
		"./testpls/asf1.asf": asf1,
		"./testpls/asf2.asf": asf2,
	}

	for filePath, rulez := range testFiles {
		t.Logf("Testing %s", filePath)
		testASFFile(filePath, rulez, t)
	}
}

func testASFFile(filePath string, rulez map[int]plTestStruct, t *testing.T) {

	rawPlaylist := getPLFile(filePath)

	parser := NewAsfParser(rawPlaylist)
	parser.Parse()

	expectedStreamCount := len(rulez) - 1
	streamCount := len(parser.Streams)

	if expectedStreamCount != streamCount {
		t.Fatalf("Expected %d streams got %d (%s)", expectedStreamCount, streamCount, filePath)
	}

	// Map of tested streams
	keys := make(map[int]bool, 10)

	for _, stream := range parser.Streams {

		expected, ok := rulez[stream.Index]

		if !ok {
			t.Fatalf("Stream with index %d was not expected", stream.Index)
		}

		keys[stream.Index] = true

		if expected.Index != stream.Index {
			t.Fatalf("Expected stream index %d == %d", stream.Index, expected.Index, stream.Index)
		}

		if expected.Title != stream.Title {
			t.Fatalf("Expected stream (%s:%d) Title '%s' == '%s'", filePath, stream.Index, expected.Title, stream.Title)
		}

		if expected.Description != stream.Description {
			t.Fatalf("Expected stream (%s:%d) Description '%s' == '%s'", filePath, stream.Index, expected.Description, stream.Description)
		}

		if expected.Logo != stream.Logo {
			t.Fatalf("Expected stream (%s:%d) Logo '%s' == '%s'", filePath, stream.Index, expected.Logo, stream.Logo)
		}

		if expected.Author != stream.Author {
			t.Fatalf("Expected stream (%s:%d) Author '%s' == '%s'", filePath, stream.Index, expected.Author, stream.Author)
		}

		if expected.Copyright != stream.Copyright {
			t.Fatalf("Expected stream (%s:%d) Copyright '%s' == '%s'", filePath, stream.Index, expected.Copyright, stream.Copyright)
		}

		if expected.MoreInfo != stream.MoreInfo {
			t.Fatalf("Expected stream (%s:%d) MoreInfo '%s' == '%s'", filePath, stream.Index, expected.MoreInfo, stream.MoreInfo)
		}

		if expected.Url != stream.Url {
			t.Fatalf("Expected stream (%s:%d) Url '%s' == '%s'", filePath, stream.Index, expected.Url, stream.Url)
		}
	}

	if len(keys) != len(rulez)-1 {
		t.Fatalf("Not all streams has been tested")
	}
}

func BenchmarkAsfParsing(b *testing.B) {

	testFile := getPLFile("./testpls/asf1.asf")

	for i := 0; i < b.N; i++ {
		parser := NewAsfParser(testFile)
		parser.Parse()
	}
}
