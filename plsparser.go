// Part of the Go playlist parser package
//
// Copyright 2013 Rafal Zajac rzajac<at>gmail<dot>com. All rights reserved.
// http://github.com/rzajac/plparser
//
// Licensed under the MIT license

// Package plparser provides primitives to parse PLS, ASX, ASF and M3U playlists.
package plparser

import (
	"bufio"
	"bytes"
	"io"
	"regexp"
	"strconv"
)

// plsRegs regular expressions to match interesting parts of PLS playlist.
var plsRegs = []struct {
	name string
	reg  *regexp.Regexp
}{
	{"Title", regexp.MustCompile(`(?is)title([0-9]+)(?:\s+)?=(?:\s+)?(.*)`)},
	{"Url", regexp.MustCompile(`(?i)file([0-9]+)(?:\s+)?=(?:\s+)?(.*)`)}}

// PlsParser implements PLS playlist parser.
type PlsParser struct {
	raw     []byte
	reader  *bufio.Reader
	Streams []*Stream
}

// NewPlsParser returns new PLS playlist parser. Takes playlist raw content to parse.
func NewPlsParser(raw []byte) *PlsParser {
	pls := new(PlsParser)
	pls.raw = raw
	pls.Streams = make([]*Stream, 0, 10)

	br := bytes.NewReader(pls.raw)
	pls.reader = bufio.NewReader(br)
	return pls
}

// Parse parses a PLS playlist.
func (p *PlsParser) Parse() {

	var titles = make(map[int]string, 10)
	var streams = make(map[int]*Stream, 10)

	for {
		line, err := p.reader.ReadString('\n')

		if err != nil && err != io.EOF {
			break
		}

		for _, s := range plsRegs {
			values := s.reg.FindStringSubmatch(line)

			if len(values) == 3 {

				v := fixString(values[2])
				idx64, _ := strconv.ParseInt(values[1], 10, 0)
				idx := int(idx64)

				if s.name == "Url" {

					stream := NewStream(idx)
					stream.Url = v

					if title, ok := titles[idx]; ok {
						stream.Title = title
					}

					streams[idx] = stream
				} else {
					if stream, ok := streams[idx]; ok {
						stream.Title = v
					} else {
						titles[idx] = v
					}
				}
			}
		}

		if err == io.EOF {
			break
		}
	}

	for _, v := range streams {
		p.Streams = append(p.Streams, v)
	}
}

// GetStreams gets list of found streams in the playlist.
func (p *PlsParser) GetStreams() []*Stream {
	return p.Streams
}
