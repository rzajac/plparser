// Part of the Go playlist parser package
//
// Copyright 2013 Rafal Zajac rzajac<at>gmail<dot>com. All rights reserved.
// http://github.com/rzajac/plparser
//
// Licensed under the MIT license

package plparser

import (
	"bufio"
	"bytes"
	// "fmt"
	"io"
	"strings"
)

func isUrl(line string) (is bool) {
	if line != "" && (strings.HasPrefix(line, "http") || strings.HasPrefix(line, "mms")) {
		is = true
	} else {
		is = false
	}
	return
}

// NewM3uParser creates new M3U playlist parser.
// Takes playlist text and returns M3U parser.
func NewM3uParser(raw []byte) (m3u *M3uParser) {
	m3u = new(M3uParser)
	m3u.raw = raw
	m3u.Streams = make([]*Stream, 0, 10)

	br := bytes.NewReader(m3u.raw)
	m3u.reader = bufio.NewReader(br)
	return
}

// M3uParser implements M3U playlist parser.
type M3uParser struct {
	raw     []byte
	reader  *bufio.Reader
	Streams []*Stream
}

// Parse parses a M3U playlist.
func (p *M3uParser) Parse() {
	var idx int

	for {
		line, err := p.reader.ReadString('\n')

		if err != nil && err != io.EOF {
			break
		}

		line = fixString(line)

		if isUrl(line) {
			idx += 1
			stream := new(Stream)
			stream.Index = idx
			stream.Url = line
			p.Streams = append(p.Streams, stream)
		}

		if err == io.EOF {
			break
		}
	}
}

// GetStreams gets list of streams in a playlist.
func (p *M3uParser) GetStreams() []*Stream {
	return p.Streams
}
