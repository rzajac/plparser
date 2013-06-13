// Part of the Go playlist parser package
//
// Copyright 2013 Rafal Zajac rzajac<at>gmail<dot>com. All rights reserved.
// http://github.com/rzajac/plparser
//
// Licensed under the MIT license

package plparser

import (
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"
)

// Binary content types.
var BINARY = map[string]bool{
	"audio/mpeg":               true,
	"application/octet-stream": true}

// Text content types.
var TEXT = map[string]bool{
	"text/plain":                true,
	"text/plain; charset=utf-8": true,
	"text/html":                 true,
	"text/html; charset=utf-8":  true,
	"audio/x-scpls":             true, // PLS playlist
	"video/x-ms-asf":            true, // ASX playlist
	"audio/mpegurl":             true, // M3U playlist
	"audio/x-mpegurl":           true, // M3U playlist
}

// HttpResp is used to send response through typed channel.
type HttpResp struct {
	resp *http.Response
	err  error
}

// getBody does HTTP GET and sends response on return channel.
func getBody(client *http.Client, req *http.Request, retch chan *HttpResp) {

	resp, err := client.Do(req)

	r := new(HttpResp)
	r.resp = resp
	r.err = err

	retch <- r
}

// PlaylistResp the playlist response.
type PlaylistResp struct {
	// Url to the playlist. For files this is a absolute path.
	Url string
	// StatusCode the HTTP response code.
	StatusCode int
	// ContentType type from HTTP response headers.
	ContentType string
	// ContentTypeDetected holds return value of http.DetectContentType().
	ContentTypeDetected string
	// Raw is the raw response. If the response was detected as binary it
	// has only first playlistReadLimit bytes.
	Raw []byte
	// Origin is where the playlist came from: ORIGIN_FILE, ORIGIN_URL
	Origin string
}

// NewPlaylistRespUrl creates new playlist response. Takes URL to potential playlist.
func NewPlaylistRespUrl(url string, timeout int) (*PlaylistResp, error) {

	plr := new(PlaylistResp)
	plr.Url = url
	plr.Origin = ORIGIN_URL
	client := &http.Client{}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return plr, err
	}

	var response *HttpResp
	var resp *http.Response

	// Channel to get response on
	retch := make(chan *HttpResp)

	go getBody(client, req, retch)

	// Wait for response or give up after timeout
	select {

	case response = <-retch:
		resp = response.resp
		err = response.err

	case <-time.After(time.Duration(timeout) * time.Second):
		resp = nil
		err = errors.New("Timeout connecting to URL.")
	}

	if resp == nil {
		// This means that ShoutCast server responded with its header
		// which is not recognized by http package.
		// The header is usually in the form of ICY 200 OK
		// In this case we set the response to be 200 but containing binary data.
		if strings.Contains(err.Error(), "malformed HTTP version \"ICY\"") {
			plr.StatusCode = 200
			plr.ContentType = "application/octet-stream"
			plr.ContentTypeDetected = "application/octet-stream"
			err = nil
		}

		return plr, err
	}

	if err != nil {
		return plr, err
	}

	defer resp.Body.Close()

	plr.StatusCode = resp.StatusCode
	plr.ContentType = resp.Header.Get("Content-Type")

	// If it's text response we read whole content
	if _, ok := TEXT[plr.ContentType]; ok {
		plr.Raw, err = ioutil.ReadAll(resp.Body)
	} else {
		// Here we read only playlistReadLimit bytes so we can use DetectContentType
		plr.Raw, err = ioutil.ReadAll(io.LimitReader(resp.Body, playlistReadLimit))

	}

	if err != nil {
		return plr, err
	}

	plr.ContentTypeDetected = http.DetectContentType(plr.Raw)

	return plr, err
}

// NewPlaylistRespFile creates new playlist response from local file.
func NewPlaylistRespFile(path string) (*PlaylistResp, error) {

	plr := new(PlaylistResp)
	plr.Url = path
	plr.Origin = ORIGIN_FILE

	file, err := os.Open(path)
	if err != nil {
		plr.StatusCode = 500
		return plr, err
	}

	plr.Raw, err = ioutil.ReadAll(file)
	if err != nil {
		plr.StatusCode = 500
		return plr, err
	}

	plr.StatusCode = 200
	plr.ContentTypeDetected = http.DetectContentType(plr.Raw)

	return plr, err
}

// IsBinary returns true if playlist content is binary.
func (pr *PlaylistResp) IsBinary() bool {
	ret := false

	if _, ok := BINARY[pr.ContentTypeDetected]; ok {
		ret = true
	}

	return ret
}

// IsHtml returns true if playlist content is HTML.
func (pr *PlaylistResp) IsHtml() bool {
	ret := false

	if pr.ContentTypeDetected == FT_HTML {
		ret = true
	}

	return ret
}

// IsPotentialPlaylist returns true if playlist content is potentially valid playlist.
func (pr *PlaylistResp) IsPotentialPlaylist() bool {
	ret := false

	if !(pr.IsBinary() || pr.IsHtml()) {
		ret = true
	}

	return ret
}
