// Part of the Go playlist parser package
//
// Copyright 2013 Rafal Zajac rzajac<at>gmail<dot>com. All rights reserved.
// http://github.com/rzajac/plparser
//
// Licensed under the MIT license

package plparser

import (
	"reflect"
)

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

	// Some unexported properties to handle parsing
	// of various playlists.
	raw  string
	Base string
}

// NewStream returns new stream. Takes stream's index in a playlist.
func NewStream(index int) *Stream {
	s := new(Stream)
	s.Index = index
	return s
}

// makeCopy makes a copy of a stream.
// NOTE: This does not copy unexported properties.
func (s *Stream) makeCopy() *Stream {

	str := new(Stream)
	str.Index = s.Index
	str.Title = s.Title
	str.Description = s.Description
	str.Logo = s.Logo
	str.Author = s.Author
	str.Copyright = s.Copyright
	str.MoreInfo = s.MoreInfo
	str.Url = s.Url

	return str
}

// setValue sets Stream structure value by name.
func (s *Stream) setValue(fieldName, value string) {
	reflect.ValueOf(s).Elem().FieldByName(fieldName).SetString(value)
}
