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
	"encoding/json"
	// "fmt"
	"io"
	"log"
	"strings"
)

// Audio formats
var (
	MPEG = []string{"audio/mpeg"}
)

// Header content types
var (
	PL_PLS = "audio/x-scpls"
)

// Remove new lines and trim
func fixString(s string) string {
	v := strings.Replace(s, "\r\n", "", -1)
	v = strings.Replace(v, "\n", "", -1)
	v = strings.Replace(v, "\r", "", -1)
	v = strings.TrimSpace(v)
	return v
}

type Playlister interface {
	Parse()
	GetStreams() []*Stream
}

// Create new playlist based on PlaylistResponse
func NewPlaylist(plr *PlaylistResp) (pl *Playlist) {

	pl = new(Playlist)
	pl.Resp = plr

	br := bytes.NewReader(pl.Resp.Raw)
	pl.lineReader = bufio.NewReader(br)

	return
}

// The playlist
type Playlist struct {
	Type    string
	Streams []*Stream
	Resp    *PlaylistResp

	firstLine  string
	lineReader *bufio.Reader
}

func (p *Playlist) Parse() string {

	// Get first line that is not empty
	fl, err := p.getLine()
	if err != nil {
		log.Fatal(err)
	}

	// We use the first line of the playlist
	// to detect playlist type
	p.firstLine = fl

	// Detect playlist and parse it
	if p.detectType() {
		var parser Playlister

		switch p.Type {
		case "pls":
			parser = NewPlsParser(p.Resp.Raw)
		case "asf":
			parser = NewAsfParser(p.Resp.Raw)
		case "asx":
			parser = NewAsxParser(p.Resp.Raw)
		case "m3u":
			parser = NewM3uParser(p.Resp.Raw)
		}

		if parser != nil {
			parser.Parse()
			p.Streams = parser.GetStreams()
		}
	}

	return p.Type
}

// Detect playlist type based on the first line
func (p *Playlist) detectType() bool {

	header := strings.ToLower(p.firstLine)

	if header == "[playlist]" {
		p.Type = "pls"
	}

	if header == "[reference]" {
		p.Type = "asf"
	}

	if strings.HasPrefix(header, "<asx") {
		p.Type = "asx"
	}

	if strings.HasPrefix(header, "http") {
		p.Type = "m3u"
	}

	if strings.HasPrefix(header, "#EXTM3U") || strings.HasPrefix(header, "#EXTINF") {
		p.Type = "m3u"
	}

	return p.IsDetected()
}

// Returns true if playlist type was detected
func (p *Playlist) IsDetected() (ret bool) {
	if p.Type == "" {
		ret = false
	} else {
		ret = true
	}

	return
}

// Get next not empty line form the playlist
func (p *Playlist) getLine() (line string, err error) {

	for {
		line, err = p.lineReader.ReadString('\n')

		if err != nil && err != io.EOF {
			return
		}

		line = fixString(line)

		if len(line) != 0 {
			break
		}

	}

	if err == io.EOF {
		err = nil
	}

	return
}

// Return streams as JSON
func (p *Playlist) StreamsAsJson() string {
	j, err := json.MarshalIndent(p.Streams, " ", " ")
	if err != nil {
		log.Fatal(err)
	}

	return string(j)
}
