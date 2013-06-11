// Part of the Go playlist parser package
//
// Copyright 2013 Rafal Zajac rzajac<at>gmail<dot>com. All rights reserved.
// http://github.com/rzajac/plparser
//
// Licensed under the MIT license

package plparser

// Stream is a struct representing a stream.
type Stream struct {
	Index       int    `json:"index"`
	Title       string `json:"title"`
	Description string `json:"descr"`
	Logo        string `json:"logo"`
	Author      string `json:"author"`
	Copyright   string `json:"copyright"`
	MoreInfo    string `json:"info"`
	Url         string `json:"url"`

	// Some private properties to handle parsing
	// of various playlists
	raw  string
	base string
}

// makeCopy makes a copy of a stream.
// NOTE: This does not copy private properties.
func (s *Stream) makeCopy() (str *Stream) {
	str = new(Stream)
	str.Index = s.Index
	str.Title = s.Title
	str.Description = s.Description
	str.Logo = s.Logo
	str.Author = s.Author
	str.Copyright = s.Copyright
	str.MoreInfo = s.MoreInfo
	str.Url = s.Url

	return
}
