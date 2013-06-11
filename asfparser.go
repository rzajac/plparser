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

// NewAsfParser creates new ASF playlist parser.
// Takes playlist text and returns ASF parser.
func NewAsfParser(raw []byte) (asf *AsfParser) {
	asf = new(AsfParser)
	asf.raw = raw
	asf.Streams = make([]*Stream, 0, 10)

	// Create a reader
	br := bytes.NewReader(asf.raw)
	asf.reader = bufio.NewReader(br)
	return
}

// AsfParser implements ASF playlist parser.
type AsfParser struct {
	raw     []byte
	reader  *bufio.Reader
	Streams []*Stream
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

			stream := new(Stream)
			stream.Url = streamUrl
			stream.Index = idx

			p.Streams = append(p.Streams, stream)
		}

		if err == io.EOF {
			break
		}
	}
}

// GetStreams gets list of streams in a playlist.
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
