// Part of the Go playlist parser package
//
// Copyright 2013 Rafal Zajac rzajac<at>gmail<dot>com. All rights reserved.
// http://github.com/rzajac/plparser
//
// Licensed under the MIT license

package plparser

import (
	"regexp"
	"strconv"
	"testing"
)

func TestAsxRegExp(t *testing.T) {

	var descreg, titlereg, logoreg, authorreg, copyreg, urlreg *regexp.Regexp
	var basereg, morereg *regexp.Regexp

	// Extract regular expressions
	for _, reg := range asxRegs {

		switch reg.name {

		case "Description":
			descreg = reg.reg

		case "Title":
			titlereg = reg.reg

		case "Logo":
			logoreg = reg.reg

		case "Author":
			authorreg = reg.reg

		case "Copyright":
			copyreg = reg.reg

		case "Url":
			urlreg = reg.reg

		case "base":
			basereg = reg.reg

		case "MoreInfo":
			morereg = reg.reg
		}
	}

	// Regular expression rules
	var testRuls = []struct {
		reg      *regexp.Regexp
		line     string
		expected string
		matches  int
	}{
		// Description (Abstract)
		{descreg, "<ABSTRACT>aaaa</ABSTRACT>", "aaaa", 1},
		{descreg, "<AbstracT>aaaa</AbStRaCt>", "aaaa", 1},
		{descreg, "<abstract>aaaa</abstract>", "aaaa", 1},
		{descreg, "<abstract>aaaa</ABSTRACT>", "aaaa", 1},
		{descreg, "<ABSTRACT>aaaa</abstract>", "aaaa", 1},
		{descreg, "<ABSTRACT  >aaaa</abstract  >", "aaaa", 1},
		{descreg, "<ABSTRACT  >aaaa</abstract>", "aaaa", 1},

		// Title
		{titlereg, "<TITLE>aaaa</TITLE>", "aaaa", 1},
		{titlereg, "<TitlE>aaaa</TiTlE>", "aaaa", 1},
		{titlereg, "<title>aaaa</title>", "aaaa", 1},
		{titlereg, "<title>aaaa</TITLE>", "aaaa", 1},
		{titlereg, "<TITLE >aaaa</TITLE >", "aaaa", 1},
		{titlereg, "<TITLE  >aaaa</TITLE>", "aaaa", 1},
		{titlereg, "<TITLE  >aaaa</TiTle>", "aaaa", 1},

		// Logo
		{logoreg, "<LOGO href='aaaa' Style='BANNER'/>", "aaaa", 1},
		{logoreg, "<LOGO Href='aaaa' Style=\"BANNER\"/>", "aaaa", 1},
		{logoreg, "<LOGO href=\"aaaa\" Style=\"BANNER\"/>", "aaaa", 1},
		{logoreg, "<LOGO Href = 'aaaa' Style='BANNER'/>", "aaaa", 1},
		{logoreg, "<LOGO  href ='aaaa' Style='BANNER'  />", "aaaa", 1},
		{logoreg, "<logo href= 'aaaa' Style='BANNER'  />", "aaaa", 1},

		// Author
		{authorreg, "<AUTHOR>aaaa</AUTHOR>", "aaaa", 1},
		{authorreg, "<author>aaaa</AUTHOR>", "aaaa", 1},
		{authorreg, "<author>aaaa</author>", "aaaa", 1},
		{authorreg, "<author >aaaa</author >", "aaaa", 1},
		{authorreg, "<author  >aaaa</author  >", "aaaa", 1},
		{authorreg, "<author>aaaa</author  >", "aaaa", 1},

		// Copyright
		{copyreg, "<COPYRIGHT>aaaa</COPYRIGHT>", "aaaa", 1},
		{copyreg, "<COPYRIGHT >aaaa</COPYRIGHT>", "aaaa", 1},
		{copyreg, "<COPYRIGHT >aaaa</COPYRIGHT >", "aaaa", 1},
		{copyreg, "<COPYRIGHT  >aaaa</COPYRIGHT  >", "aaaa", 1},
		{copyreg, "<copyright  >aaaa</COPYRIGHT  >", "aaaa", 1},
		{copyreg, "<copyright  >aaaa</copyright  >", "aaaa", 1},
		{copyreg, "<Copyright  >aaaa</Copyright  >", "aaaa", 1},

		// Url (REF)
		{urlreg, "<ref href=\"aaaa\"/>", "aaaa", 1},
		{urlreg, "<ref href= \"aaaa\"/>", "aaaa", 1},
		{urlreg, "<ref href = \"aaaa\"/>", "aaaa", 1},
		{urlreg, "<ref href =\"aaaa\"/>", "aaaa", 1},
		{urlreg, "<ref Href='aaaa'/>", "aaaa", 1},
		{urlreg, "<ref href= 'aaaa'/>", "aaaa", 1},
		{urlreg, "<ref href = 'aaaa'/>", "aaaa", 1},
		{urlreg, "<ref href ='aaaa'/>", "aaaa", 1},
		{urlreg, "<REF Href ='aaaa'/>", "aaaa", 1},
		{urlreg, "<Ref href ='aaaa'/>", "aaaa", 1},
		{urlreg, "<ReF href ='aaaa'/>", "aaaa", 1},
		// Closed REF
		{urlreg, "<ref href=\"aaaa\"></ref>", "aaaa", 1},
		{urlreg, "<ref href= \"aaaa\"></ref>", "aaaa", 1},
		{urlreg, "<ref href = \"aaaa\"></ref>", "aaaa", 1},
		{urlreg, "<ref href =\"aaaa\"></ref>", "aaaa", 1},
		{urlreg, "<ref href='aaaa'></ref>", "aaaa", 1},
		{urlreg, "<ref href= 'aaaa'></ref>", "aaaa", 1},
		{urlreg, "<ref Href = 'aaaa'></ref>", "aaaa", 1},
		{urlreg, "<ref href ='aaaa'></ref>", "aaaa", 1},
		{urlreg, "<REF href ='aaaa'></ref>", "aaaa", 1},
		{urlreg, "<Ref href ='aaaa'></ref >", "aaaa", 1},
		{urlreg, "<ReF Href ='aaaa'></ref >", "aaaa", 1},

		// Base
		{basereg, "<BASE href=\"aaaa\"/>", "aaaa", 1},
		{basereg, "<BasE href= \"aaaa\"/>", "aaaa", 1},
		{basereg, "<base Href = \"aaaa\" />", "aaaa", 1},
		{basereg, "<Base href =\"aaaa\"  />", "aaaa", 1},

		{basereg, "<BASE href='aaaa'/>", "aaaa", 1},
		{basereg, "<BasE Href= 'aaaa'/>", "aaaa", 1},
		{basereg, "<base href = 'aaaa' />", "aaaa", 1},
		{basereg, "<Base href ='aaaa'  />", "aaaa", 1},

		// Closed BASE
		{basereg, "<BASE href=\"aaaa\"></BASE>", "aaaa", 1},
		{basereg, "<BasE Href= \"aaaa\" ></Base>", "aaaa", 1},
		{basereg, "<base href = \"aaaa\" ></BASE>", "aaaa", 1},
		{basereg, "<Base href =\"aaaa\" ></basE>", "aaaa", 1},

		{basereg, "<BASE href='aaaa'></BASE>", "aaaa", 1},
		{basereg, "<BasE href= 'aaaa'></BasE>", "aaaa", 1},
		{basereg, "<base href = 'aaaa'></baSe>", "aaaa", 1},
		{basereg, "<Base Href ='aaaa'> </base>", "aaaa", 1},

		// MoreInfo
		{morereg, "<MOREINFO href=\"aaaa\" />", "aaaa", 1},
		{morereg, "<MoreinfO href =\"aaaa\" />", "aaaa", 1},
		{morereg, "<moreinfo Href = \"aaaa\"/>", "aaaa", 1},
		{morereg, "<MOREINFO href= \"aaaa\"  />", "aaaa", 1},

		// Closed MoreInfo
		{morereg, "<MOREINFO href=\"aaaa\" > </MOREINFO>", "aaaa", 1},
		{morereg, "<MOREINFO href =\"aaaa\"></MOREINFO>", "aaaa", 1},
		{morereg, "<moreinfo href = \"aaaa\" ></moreinfo>", "aaaa", 1},
		{morereg, "<MOREINFO href= \"aaaa\" ></MoreinfO>", "aaaa", 1},
	}

	for _, rule := range testRuls {

		values := rule.reg.FindStringSubmatch(rule.line)

		expectedMatches := rule.matches + 1
		matchesFound := len(values)

		if expectedMatches != matchesFound {
			t.Fatalf("Expected matches %d == %d (%s)", expectedMatches, matchesFound, rule.line)
		}

		if rule.expected != values[1] {
			t.Fatalf("Expected regexp to return '%s' but got '%s' for line '%s'", rule.expected, values[1], rule.line)
		}

	}
}

func TestAsxFiles(t *testing.T) {

	var asx1 = map[int][]plTestStruct{
		// Entry number
		1: {
			// Streams in the entry
			{}, // Streams in entries are 1 indexed
			{1, "E1T", "E1A", "http://E1.ex.com/l.gif", "E1AU", "E1C", "http://E1.mi.ex.com", "http://E1.st1.ex.com:8881/"},
			{2, "E1T", "E1A", "http://E1.ex.com/l.gif", "E1AU", "E1C", "http://E1.mi.ex.com", "http://E1.st2.ex.com:8882/"},
		},
		2: {
			{},
			{1, "E2T", "E2A", "http://E2.ex.com/l.gif", "E2AU", "E2C", "http://E2.mi.ex.com", "http://E2.st1.ex.com:8881/"},
			{2, "E2T", "E2A", "http://E2.ex.com/l.gif", "E2AU", "E2C", "http://E2.mi.ex.com", "http://E2.st2.ex.com:8882/"},
		},
	}

	var asx2 = map[int][]plTestStruct{
		1: {
			{},
			{1, "MT", "MA", "http://ml.ex.com/l.gif", "MAU", "MC", "http://mi.ex.com/mi", "http://E1.st1.ex.com:8881/"},
			{2, "MT", "MA", "http://ml.ex.com/l.gif", "MAU", "MC", "http://mi.ex.com/mi", "http://E1.st2.ex.com:8882/"},
		},
		2: {
			{},
			{1, "MT", "E2A", "http://ml.ex.com/l.gif", "MAU", "MC", "http://E2.mi.ex.com", "http://E2.st1.ex.com:8882/"},
		},
	}

	// Files to test
	var testFiles = map[string]map[int][]plTestStruct{
		"./testpls/asx1.asx": asx1,
		"./testpls/asx2.asx": asx2,
	}

	// Test each file
	for filePath, rulez := range testFiles {
		testAsxFile(filePath, rulez, t)
	}
}

func testAsxFile(filePath string, rulez map[int][]plTestStruct, t *testing.T) {

	testFile := getPLFile(filePath)
	parser := NewAsxParser(testFile)
	parser.Parse()

	var entry []plTestStruct
	var ok bool

	// Map of tested streams in form X-Y (1-1, 1-2)
	// Where X is the entry number and Y is stream number
	keys := make(map[string]bool, 10)

	// Go over extracted streams and check if data in them matches
	// the tests
	for _, stream := range parser.Streams {

		// First we get entry and stream number from the REF URL
		// We do that to know from which ENTRY this stream comes from
		// and which stream are we testing. In the test files we this data
		// in the stream URL, here we just extract it
		entryNumb, streamNumber, key := getStreamNumbers(stream.Url, t)
		keys[key] = true

		// t.Logf("en:%d, sn:%d", entryNumb, streamNumber)

		if entry, ok = rulez[entryNumb]; !ok {
			t.Fatalf("Entry with index %d was not expected", entryNumb)
		}

		streamExpected := entry[streamNumber]

		if streamExpected.Title != stream.Title {
			t.Fatalf("Expected stream (%s:%d:%d) Title '%s' == '%s'", filePath, entryNumb, streamNumber, streamExpected.Title, stream.Title)
		}

		if streamExpected.Description != stream.Description {
			t.Fatalf("Expected stream (%s:%d:%d) Description '%s' == '%s'", filePath, entryNumb, streamNumber, streamExpected.Description, stream.Description)
		}

		if streamExpected.Logo != stream.Logo {
			t.Fatalf("Expected stream (%s:%d:%d) Logo '%s' == '%s'", filePath, entryNumb, streamNumber, streamExpected.Logo, stream.Logo)
		}

		if streamExpected.Author != stream.Author {
			t.Fatalf("Expected stream (%s:%d:%d) Author '%s' == '%s'", filePath, entryNumb, streamNumber, streamExpected.Author, stream.Author)
		}

		if streamExpected.Copyright != stream.Copyright {
			t.Fatalf("Expected stream (%s:%d:%d) Copyright '%s' == '%s'", filePath, entryNumb, streamNumber, streamExpected.Copyright, stream.Copyright)
		}

		if streamExpected.MoreInfo != stream.MoreInfo {
			t.Fatalf("Expected stream (%s:%d:%d) MoreInfo '%s' == '%s'", filePath, entryNumb, streamNumber, streamExpected.MoreInfo, stream.MoreInfo)
		}

		if streamExpected.Url != stream.Url {
			t.Fatalf("Expected stream (%s:%d:%d) Url '%s' == '%s'", filePath, entryNumb, streamNumber, streamExpected.Url, stream.Url)
		}
	}

	var expectedStreamCount int

	// Count number of expected streams in rulez and
	// check if all expected streams were tested.
	for entrynumb, strrule := range rulez {
		key := strconv.FormatInt(int64(entrynumb), 10) + "-"
		for _, str := range strrule {

			if str.Index > 0 {
				tmpKey := key + strconv.FormatInt(int64(str.Index), 10)

				if _, ok = keys[tmpKey]; !ok {
					t.Fatalf("Stream %s not present in the result", tmpKey)
				}
				expectedStreamCount += 1
			}
		}
	}

	if expectedStreamCount != len(parser.Streams) {
		t.Fatalf("Expected %d streams got %d (%s)", expectedStreamCount, len(parser.Streams), filePath)
	}
}

// Get entry number and stream number
func getStreamNumbers(url string, t *testing.T) (enumb, stnumb int, key string) {

	regRef := regexp.MustCompile("http://E([0-9]+)\\.st([0-9]+)\\.ex.com:888")
	values := regRef.FindStringSubmatch(url)

	if len(values) == 3 {
		key = values[1] + "-" + values[2]
		i64, _ := strconv.ParseInt(values[1], 10, 0)
		enumb = int(i64)
		i64, _ = strconv.ParseInt(values[2], 10, 0)
		stnumb = int(i64)
	} else {
		panic("Can't get entry number")
	}

	return
}

func BenchmarkAsxParsing(b *testing.B) {

	testFile := getPLFile("./testpls/asx2.asx")

	for i := 0; i < b.N; i++ {
		parser := NewAsxParser(testFile)
		parser.Parse()
	}
}

func BenchmarkAsxParsing2(b *testing.B) {

	testFile := getPLFile("./testpls/asx1.asx")

	for i := 0; i < b.N; i++ {
		parser := NewAsxParser(testFile)
		parser.Parse()
	}
}
