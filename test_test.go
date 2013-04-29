// Part of the Go playlist parser package
//
// Copyright 2013 Rafal Zajac rzajac<at>gmail<dot>com. All rights reserved.
// http://github.com/rzajac/plparser
//
// Licensed under the MIT license

package plparser

import (
	"io/ioutil"
	"os"
)

type plTestStruct struct {
	Index       int
	Title       string
	Description string
	Logo        string
	Author      string
	Copyright   string
	MoreInfo    string
	Url         string
}

func getPLFile(path string) (rawFile []byte) {

	file, err := os.Open(path)
	if err != nil {
		panic("Can not open test file: " + path)
	}

	rawFile, err = ioutil.ReadAll(file)
	if err != nil {
		panic(err.Error())
	}

	return
}
