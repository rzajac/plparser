// Part of the Go playlist parser package
//
// Copyright 2013 Rafal Zajac rzajac<at>gmail<dot>com. All rights reserved.
// http://github.com/rzajac/plparser
//
// Licensed under the MIT license

package plparser

import (
	"reflect"
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
	{"Base", regexp.MustCompile(`(?i)<base(?:\s+)?href(?:\s+)?=(?:\s+)?(?:"|')(.*?)(?:"|')(?:.*?)/?>(?:</base(?:\s+)?>)?`)},
	{"MoreInfo", regexp.MustCompile(`(?i)<moreinfo(?:\s+)?href(?:\s+)?=(?:\s+)?(?:"|')(.*?)(?:"|')(?:.*?)/?>(?:</moreinfo(?:\s+)?>)?`)},
}

// asxEntityRegExp regular expression to find all ENTRY elements.
var asxEntityRegExp *regexp.Regexp = regexp.MustCompile(`(?is)<entry(?:\s+)?>(.*?)</entry(?:\s+)?>`)

// AsxParser implements ASX playlist parser.
type AsxParser struct {
	raw         string
	Author      string
	Base        string
	Copyright   string
	Description string
	Logo        string
	MoreInfo    string
	Streams     []*Stream
	Title       string
}

// NewAsxParser returns new ASX playlist parser. Takes playlist raw content to parse.
func NewAsxParser(raw []byte) *AsxParser {
	asx := new(AsxParser)
	asx.raw = string(raw)
	asx.Streams = make([]*Stream, 0, 10)
	return asx
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

		// We skip URLs.
		// We first have to parse main body.
		if s.name == "Url" {
			continue
		}

		values := s.reg.FindStringSubmatch(a.raw)
		if len(values) != 2 {
			continue
		}

		a.setValue(s.name, values[1])
	}

	// Main body of the playlist has been parsed.
	// We parsed main body first to get BASE value if it exists.
	a.parseEntries(entries)
}

// GetStreams gets list of streams found in the playlist.
func (p *AsxParser) GetStreams() []*Stream {
	return p.Streams
}

// setValue sets AsxParser structure value by name.
func (a *AsxParser) setValue(fieldName, value string) {
	reflect.ValueOf(a).Elem().FieldByName(fieldName).SetString(value)
}

// getValue gets AsxParser structure field by name.
func (a *AsxParser) getValue(fieldName string) string {
	return reflect.ValueOf(a).Elem().FieldByName(fieldName).String()
}

// parseAsx parses ENTRY nodes in the ASX playlist.
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

		var value string

		if len(values) == 2 {
			value = values[1]
		} else {
			value = asxp.getValue(reg.name)
		}

		s.setValue(reg.name, value)
	}

	// Inherit base for URLs from main playlist body
	if s.Base != "" {
		s.Base = asxp.Base
	}

	// Make sure the base ends with "/"
	if s.Base != "" && !strings.HasSuffix(s.Base, "/") {
		s.Base += "/"
	}

	// Find all the stream URLs
	streams := urlr.FindAllStringSubmatch(s.raw, -1)

	var streamsToAdd []*Stream = make([]*Stream, 0, 10)

	for _, stream := range streams {

		newStream := s.makeCopy()

		// Prefix base URL to the stream URL
		if newStream.Base != "" {
			stream[1] = newStream.Base + stream[1]
		}

		newStream.Url = stream[1]
		streamsToAdd = append(streamsToAdd, newStream)
	}

	return streamsToAdd
}

// parseEntries parses ENTRY elements.
func (a *AsxParser) parseEntries(entries [][]string) {

	// Regular expression to match stream URL
	var urlr *regexp.Regexp
	var s *Stream

	// Go over all found entries
	for idx, entry := range entries {

		s = NewStream(idx)

		// First we parse all the info we can get except the URL to a stream
		for _, reg := range asxRegs {

			// We do this one last
			if reg.name == "Url" {
				urlr = reg.reg
				continue
			}

			values := reg.reg.FindStringSubmatch(entry[1])

			var value string

			if len(values) == 2 {
				value = values[1]
			} else {
				value = a.getValue(reg.name)
			}

			s.setValue(reg.name, value)
		}

		// Inherit base for URLs from main playlist body
		if s.Base != "" {
			s.Base = a.Base
		}

		// Make sure the base ends with "/"
		if s.Base != "" && !strings.HasSuffix(s.Base, "/") {
			s.Base += "/"
		}

		// Find all the stream URLs
		streams := urlr.FindAllStringSubmatch(entry[1], -1)

		for _, stream := range streams {

			newStream := s.makeCopy()

			// Prefix base URL to the stream URL
			if newStream.Base != "" {
				stream[1] = newStream.Base + stream[1]
			}

			newStream.Url = stream[1]
			a.Streams = append(a.Streams, newStream)
		}
	}
}
