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
	// "fmt"
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

// NewPlsParser creates new PLS playlist parser.
// Takes playlist text and returns PLS parser.
func NewPlsParser(raw []byte) (pls *PlsParser) {
	pls = new(PlsParser)
	pls.raw = raw
	pls.Streams = make([]*Stream, 0, 10)

	br := bytes.NewReader(pls.raw)
	pls.reader = bufio.NewReader(br)
	return
}

// PlsParser is the PLS parser.
type PlsParser struct {
	raw     []byte
	reader  *bufio.Reader
	Streams []*Stream
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

					stream := new(Stream)
					stream.Url = v
					stream.Index = idx

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

// GetStreams gets list of streams in a playlist.
func (p *PlsParser) GetStreams() []*Stream {
	return p.Streams
}
