// Part of the Go playlist parser package
//
// Copyright 2013 Rafal Zajac rzajac<at>gmail<dot>com. All rights reserved.
// http://github.com/rzajac/plparser
//
// Licensed under the MIT license

package plparser

import (
	"regexp"
	"testing"
)

func TestPlsRegExp(t *testing.T) {

	var urlreg, titlereg *regexp.Regexp

	var testUrlLines = []struct {
		line     string
		expected string
		idx      string
		matches  int
	}{
		{"File1=http://live.example.com:8888", "http://live.example.com:8888", "1", 3},
		{"file2=http://live.example.com:8888", "http://live.example.com:8888", "2", 3},
		{"file1 = http://live.example.com:8888", "http://live.example.com:8888", "1", 3},
		{"FILE2 =http://live.example.com:8888", "http://live.example.com:8888", "2", 3},
		{"FiLe2= http://live.example.com:8888", "http://live.example.com:8888", "2", 3},
		{"", "", "0", 0},
		{"not valid line", "", "0", 0},
	}

	var testTitleLines = []struct {
		line     string
		expected string
		idx      string
		matches  int
	}{
		{"Title1=Some title text", "Some title text", "1", 3},
		{"title1=Some title text", "Some title text", "1", 3},
		{"TitLe1=Some title text", "Some title text", "1", 3},
		{"TitLe1 = Some title text", "Some title text", "1", 3},
		{"TitLe1 =Some title text", "Some title text", "1", 3},
		{"TITLE1= Some title text", "Some title text", "1", 3},
		{"", "", "0", 0},
		{"not valid line", "", "0", 0},
	}

	for _, reg := range plsRegs {
		switch reg.name {

		case "Title":
			titlereg = reg.reg
		case "Url":
			urlreg = reg.reg
		}
	}

	for _, testLine := range testUrlLines {
		values := urlreg.FindStringSubmatch(testLine.line)

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

	for _, testLine := range testTitleLines {
		values := titlereg.FindStringSubmatch(testLine.line)

		if len(values) != testLine.matches {
			t.Fatalf("Expected to find %d matches for %s", testLine.matches, testLine.line)
		}

		if testLine.matches > 0 {

			if values[2] != testLine.expected {
				t.Fatalf("Expected title %s == %s", testLine.expected, values[2])
			}

			if values[1] != testLine.idx {
				t.Fatalf("Expected idx %s == %s", testLine.expected, values[1])
			}
		}
	}

}

func TestPlsFiles(t *testing.T) {

	var pls1 = map[int]plTestStruct{
		0: {},
		1: {1, "", "", "", "", "", "", "http://live.example.com:8881/"},
		2: {2, "", "", "", "", "", "", "http://live.example.com:8882/"},
		3: {3, "", "", "", "", "", "", "http://live.example.com:8883/"},
		4: {4, "", "", "", "", "", "", "http://live.example.com:8884/"},
		5: {5, "Some example title", "", "", "", "", "", "http://live.example.com:8885/"},
	}

	var pls2 = map[int]plTestStruct{
		0: {},
		1: {1, "", "", "", "", "", "", "http://live.example.com:8881/"},
		2: {2, "", "", "", "", "", "", "http://live.example.com:8882/"},
		3: {3, "", "", "", "", "", "", "http://live.example.com:8883/"},
		4: {4, "Some example title4", "", "", "", "", "", "http://live.example.com:8884/"},
		5: {5, "Some example title5", "", "", "", "", "", "http://live.example.com:8885/"},
	}

	var testFiles = map[string]map[int]plTestStruct{
		"./testpls/pls1.pls": pls1,
		"./testpls/pls2.pls": pls2,
		"./testpls/pls3.pls": pls2,
		"./testpls/pls4.pls": pls2,
	}

	for filePath, rulez := range testFiles {
		testPLSFile(filePath, rulez, t)
	}
}

func testPLSFile(filePath string, rulez map[int]plTestStruct, t *testing.T) {

	rawPlaylist := getPLFile(filePath)

	parser := NewPlsParser(rawPlaylist)
	parser.Parse()

	streamCount := len(parser.Streams)
	expectedStreamCount := len(rulez) - 1

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

func BenchmarkPlsParsing(b *testing.B) {

	testFile := getPLFile("./testpls/pls1.pls")

	for i := 0; i < b.N; i++ {
		parser := NewPlsParser(testFile)
		parser.Parse()
	}
}
