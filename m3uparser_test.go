package plparser

import (
	"testing"
)

func TestM3uRegExp(t *testing.T) {

	var testUrlLines = []struct {
		line  string
		isUrl bool
	}{
		{"http://live.str3am.com:2150/", true},
		{"mms://live.str3am.com:2150/", true},
		{"", false},
		{"not valid line", false},
	}

	for _, testLine := range testUrlLines {

		if isUrl(testLine.line) != testLine.isUrl {
			t.Fatalf("Expected isUrl('%s') to %v", testLine.line, testLine.isUrl)

		}
	}
}

func TestM3uFiles(t *testing.T) {

	var m3u1 = map[int]plTestStruct{
		0: {},
		1: {1, "", "", "", "", "", "", "http://live1.example.com:2151/"},
		2: {2, "", "", "", "", "", "", "http://live2.example.com:2152/"},
	}

	var m3u2 = map[int]plTestStruct{
		0: {},
		1: {1, "", "", "", "", "", "", "http://live1.example.com:2151/"},
		2: {2, "", "", "", "", "", "", "http://live2.example.com:2152/"},
		3: {3, "", "", "", "", "", "", "http://live3.example.com:2153/"},
	}

	var testFiles = map[string]map[int]plTestStruct{
		"./testpls/m3u1.m3u": m3u1,
		"./testpls/m3u2.m3u": m3u2,
	}

	for filePath, rulez := range testFiles {
		t.Logf("Testing %s", filePath)
		testM3uFile(filePath, rulez, t)
	}
}

func testM3uFile(filePath string, rulez map[int]plTestStruct, t *testing.T) {

	rawPlaylist := getPLFile(filePath)

	parser := NewM3uParser(rawPlaylist)
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

func BenchmarkM3uParsing(b *testing.B) {

	testFile := getPLFile("./testpls/m3u1.m3u")

	for i := 0; i < b.N; i++ {
		parser := NewM3uParser(testFile)
		parser.Parse()
	}
}
