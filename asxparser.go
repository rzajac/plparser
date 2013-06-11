// Part of the Go playlist parser package
//
// Copyright 2013 Rafal Zajac rzajac<at>gmail<dot>com. All rights reserved.
// http://github.com/rzajac/plparser
//
// Licensed under the MIT license

package plparser

import (
	// "fmt"
	"regexp"
	"strings"
)

// asxRegs are regular expressions to catch various ASX playlist elements.
var asxRegs = []struct {
	name string
	reg  *regexp.Regexp
}{
	{"Description", regexp.MustCompile(`(?is)<abstract(?:\s+)?>(.*?)</abstract(?:\s+)?>`)},
	{"Title", regexp.MustCompile(`(?is)<title(?:\s+)?>(.*?)</title(?:\s+)?>`)},
	{"Logo", regexp.MustCompile(`(?i)<logo(?:\s+)?href(?:\s+)?=(?:\s+)?(?:"|')(.*?)(?:"|')(?:.*?)/>`)},
	{"Author", regexp.MustCompile(`(?is)<author(?:\s+)?>(.*?)</author(?:\s+)?>`)},
	{"Copyright", regexp.MustCompile(`(?is)<copyright(?:\s+)?>(.*?)</copyright(?:\s+)?>`)},
	{"Url", regexp.MustCompile(`(?i)<ref(?:\s+)?href(?:\s+)?=(?:\s+)?(?:"|')(.*?)(?:"|')(?:.*?)/?>(?:</ref(?:\s+)?>)?`)},
	{"base", regexp.MustCompile(`(?i)<base(?:\s+)?href(?:\s+)?=(?:\s+)?(?:"|')(.*?)(?:"|')(?:.*?)/?>(?:</base(?:\s+)?>)?`)},
	{"MoreInfo", regexp.MustCompile(`(?i)<moreinfo(?:\s+)?href(?:\s+)?=(?:\s+)?(?:"|')(.*?)(?:"|')(?:.*?)/?>(?:</moreinfo(?:\s+)?>)?`)},
}

// asxEntityRegExp regular expression to find all ENTRY elements.
var asxEntityRegExp *regexp.Regexp = regexp.MustCompile(`(?is)<entry(?:\s+)?>(.*?)</entry(?:\s+)?>`)

// NewAsxParser creates new ASX playlist parser.
// Takes playlist text and returns ASX parser.
func NewAsxParser(raw []byte) (asx *AsxParser) {
	asx = new(AsxParser)
	asx.raw = string(raw)
	asx.Streams = make([]*Stream, 0, 10)
	return
}

// AsxParser implements ASX playlist parser.
type AsxParser struct {
	raw         string
	Author      string
	base        string
	Copyright   string
	Description string
	Logo        string
	MoreInfo    string
	Streams     []*Stream
	Title       string
}

// Parse parses an ASX playlist.
func (a *AsxParser) Parse() {

	// Get all the entries that represent streams
	entries := asxEntityRegExp.FindAllStringSubmatch(a.raw, -1)

	// Remove parser entries from ASX
	// This will simplify parsing the main body of the
	// playlist
	a.raw = asxEntityRegExp.ReplaceAllString(a.raw, "")

	// Main body of the playlist may have the same
	// element names as ENTRY except REF (Url)
	for _, s := range asxRegs {

		if s.name == "Url" {
			continue
		}

		values := s.reg.FindStringSubmatch(a.raw)
		valuesCount := len(values)

		switch s.name {
		case "Description":
			if valuesCount == 2 {
				a.Description = values[1]
			}

		case "Title":
			if valuesCount == 2 {
				a.Title = values[1]
			}

		case "Logo":
			if valuesCount == 2 {
				a.Logo = values[1]
			}

		case "Author":
			if valuesCount == 2 {
				a.Author = values[1]
			}

		case "Copyright":
			if valuesCount == 2 {
				a.Copyright = values[1]
			}

		case "base":
			if valuesCount == 2 {
				a.base = values[1]
			}

		case "MoreInfo":
			if valuesCount == 2 {
				a.MoreInfo = values[1]
			}
		}
	}

	// Main body of the playlist has been parsed now we
	// can parse entries we parsed main body first
	// to get BASE value if it exists
	for i, e := range entries {

		str := new(Stream)
		str.Index = i
		str.raw = e[1]

		a.Streams = append(a.Streams, str.parseAsx(a)...)
	}
}

// GetStreams gets list of streams in a playlist.
func (p *AsxParser) GetStreams() []*Stream {
	return p.Streams
}

// parseAsx makes use of Stream structure to parse entries in the ASX playlist.
func (s *Stream) parseAsx(asxp *AsxParser) []*Stream {

	// Regular expression to match stream URL
	var urlr *regexp.Regexp

	// First we parse all the info we can get except
	// the URL to a stream
	for _, reg := range asxRegs {

		// We do this one last
		if reg.name == "Url" {
			urlr = reg.reg
			continue
		}

		values := reg.reg.FindStringSubmatch(s.raw)
		valuesCount := len(values)

		switch reg.name {
		case "Description":
			if valuesCount == 2 {
				s.Description = values[1]
			} else {
				s.Description = asxp.Description
			}

		case "Title":
			if valuesCount == 2 {
				s.Title = values[1]
			} else {
				s.Title = asxp.Title
			}

		case "Logo":
			if valuesCount == 2 {
				s.Logo = values[1]
			} else {
				s.Logo = asxp.Logo
			}

		case "Author":
			if valuesCount == 2 {
				s.Author = values[1]
			} else {
				s.Author = asxp.Author
			}

		case "Copyright":
			if valuesCount == 2 {
				s.Copyright = values[1]
			} else {
				s.Copyright = asxp.Copyright
			}

		case "base":
			if valuesCount == 2 {
				s.base = values[1]
			} else {
				s.base = asxp.base
			}

		case "MoreInfo":
			if valuesCount == 2 {
				s.MoreInfo = values[1]
			} else {
				s.MoreInfo = asxp.MoreInfo
			}
		}
	}

	// Inherit base for URLs from main playlist body
	if s.base != "" {
		s.base = asxp.base
	}

	// Make sure the base ends with "/"
	if s.base != "" && !strings.HasSuffix(s.base, "/") {
		s.base += "/"
	}

	// Find all the stream URLs
	streams := urlr.FindAllStringSubmatch(s.raw, -1)

	var streamsToAdd []*Stream = make([]*Stream, 0, 10)

	for _, stream := range streams {

		newStream := s.makeCopy()

		// Prefix base URL to the stream URL
		if newStream.base != "" {
			stream[1] = newStream.base + stream[1]
		}

		newStream.Url = stream[1]
		streamsToAdd = append(streamsToAdd, newStream)
	}

	return streamsToAdd
}
