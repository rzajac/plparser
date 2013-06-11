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

// fixString removes new lines and trims a string.
func fixString(s string) string {
	v := strings.Replace(s, "\r\n", "", -1)
	v = strings.Replace(v, "\n", "", -1)
	v = strings.Replace(v, "\r", "", -1)
	v = strings.TrimSpace(v)
	return v
}

// Playlister is an interface all playlist parsers must implement.
type Playlister interface {
	// Parse parses a playlist.
	Parse()
	// GetStreams gets list of streams in a playlist.
	GetStreams() []*Stream
}

// NewPlaylist creates new playlist based on PlaylistResponse.
func NewPlaylist(plr *PlaylistResp) *Playlist {

	pl := new(Playlist)
	pl.Resp = plr

	br := bytes.NewReader(pl.Resp.Raw)
	pl.lineReader = bufio.NewReader(br)

	return pl
}

// Playlist the playlist.
type Playlist struct {
	Type    string        `json:"type"`
	Streams []*Stream     `json:"streams"`
	Resp    *PlaylistResp `json:"-"`

	firstLine  string        `json:"-"`
	lineReader *bufio.Reader `json:"-"`
}

func (p *Playlist) Parse() (string, error) {

	var err error

	// Get first line that is not empty
	// We use the first not empty line of the playlist
	// to detect playlist type
	p.firstLine, err = p.getLine()
	if err != nil {
		return p.Type, err
	}

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

	return p.Type, err
}

// detectType detects playlist type based on its first not empty line.
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

// IsDetected returns true if playlist type was detected.
func (p *Playlist) IsDetected() bool {
	ret := false

	if p.Type != "" {
		ret = true
	}

	return ret
}

// getLine gets next not empty line form a playlist.
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

// StreamsAsJson returns streams as JSON.
func (p *Playlist) StreamsAsJson() (string, error) {
	j, err := json.MarshalIndent(p, " ", " ")
	if err != nil {
		return "", NewPlParserError(err.Error(), true)
	}

	return string(j), err
}

// PlParserError 
type PlParserError struct {
	doJsonError bool
	Msg         string `json:"error"`
}

// NewPlParserError creates new error.
func NewPlParserError(msg string, doJsonError bool) *PlParserError {
	plpe := new(PlParserError)
	plpe.doJsonError = doJsonError
	plpe.Msg = msg

	return plpe
}

// Error returns error message.
func (p *PlParserError) Error() (msg string) {

	if p.doJsonError {
		e, _ := json.MarshalIndent(p, " ", " ")
		msg = string(e)
	} else {
		msg = p.Msg
	}

	return
}
