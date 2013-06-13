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
	"io"
	"regexp"
	"strconv"
)

// asfReg is a regular expression to match streams in a ASF playlist.
var asfReg *regexp.Regexp = regexp.MustCompile(`(?i)ref([0-9]+)(?:\s+)?=(?:\s+)?(.*)`)

// AsfParser implements ASF playlist parser.
type AsfParser struct {
	raw     []byte        // Raw contents of a playlist
	reader  *bufio.Reader //
	Streams []*Stream     // The array of found strams
}

// NewAsfParser returns new ASF playlist parser. Takes playlist raw content to parse.
func NewAsfParser(raw []byte) *AsfParser {
	asf := new(AsfParser)
	asf.raw = raw
	asf.Streams = make([]*Stream, 0, 10)

	// Create a reader
	br := bytes.NewReader(asf.raw)
	asf.reader = bufio.NewReader(br)
	return asf
}

// Parse parses an ASF playlist.
func (p *AsfParser) Parse() {

	for {
		line, err := p.reader.ReadString('\n')

		if err != nil && err != io.EOF {
			break
		}

		idx, streamUrl := findMatch(line, asfReg)

		if streamUrl != "" {

			stream := NewStream(idx)
			stream.Url = streamUrl

			p.Streams = append(p.Streams, stream)
		}

		if err == io.EOF {
			break
		}
	}
}

// GetStreams gets list of streams found in the playlist.
func (p *AsfParser) GetStreams() []*Stream {
	return p.Streams
}

// findMatch checks if the line matches the regular expression for stream URL.
func findMatch(line string, reg *regexp.Regexp) (idx int, stream string) {

	value := reg.FindStringSubmatch(line)

	if len(value) == 3 {

		stream = fixString(value[2])
		idx64, _ := strconv.ParseInt(value[1], 10, 0)
		idx = int(idx64)
	}

	return
}
