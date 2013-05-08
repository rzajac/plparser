# Playlist parser in Go

Playlist parser / detector in pure Go.


# Supported playlist formats

* PLS
* ASX
* ASF
* M3U

# Installation

Run

    go get github.com/rzajac/plparser

Include in your source:

    import "github.com/rzajac/plparser"

# Updating

Run

    go get -u github.com/rzajac/plparser

# Benchmarking

    go test -bench=".*" github.com/rzajac/plparser

# Usage

	// Get playlist from URL with timeout set to 5 seconds
    plr, err := plparser.NewPlaylistRespUrl("http://example.com/some_playlist", 5)

	// or get playlist from file
	plr, err := plparser.NewPlaylistRespFile("/path/to/playlist")

	if err := nil {
		// Handle error
	}

	if !(plr.StatusCode >= 200 && plr.StatusCode < 300) {
		// Handle error
	}

	if !plr.IsPotentialPlaylist() {
		// File or URL is not a playlist
	} else {
		pl := plparser.NewPlaylist(plr)
		pl.Parse()

		// If everything goes well
		// pl.Streams will have a list of streams
		// see stream.go
	}

# TODO

* Write rests for plresp.go
* Simplify tests
* Write more tests
* Add more playlist types

# License

Licensed under the MIT license

