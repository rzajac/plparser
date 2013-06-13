// Part of the Go playlist parser package
//
// Copyright 2013 Rafal Zajac rzajac<at>gmail<dot>com. All rights reserved.
// http://github.com/rzajac/plparser
//
// Licensed under the MIT license

package plparser

import (
	"strings"
)

// fixString removes new lines and trims a string.
func fixString(s string) string {
	v := strings.Replace(s, "\r\n", "", -1)
	v = strings.Replace(v, "\n", "", -1)
	v = strings.Replace(v, "\r", "", -1)
	v = strings.TrimSpace(v)
	return v
}

// isUrl is a helper function returning true if passed text is a URL.
func isUrl(text string) bool {

	var ok bool

	if text != "" && (strings.HasPrefix(text, "http") || strings.HasPrefix(text, "mms")) {
		ok = true
	}

	return ok
}
