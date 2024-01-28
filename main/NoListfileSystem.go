package main

import (
	"net/http"
	"path/filepath"
	"strings"
)

type noListFileSystem struct {
	fs http.FileSystem
}

// Open() function for noListFileSystem struct that wraps the normal fs.Open functionality.
// Inhibits the default FileServer "if it's a directory without index.html,
// show a directory listing" functionality.  With this, it will return a file not found
// error, which results in a client 404 NOT FOUND rather than offering dir listing.
func (nfs noListFileSystem) Open(path string) (http.File, error) {

	// If neither existing file nor directory, just return error so http returns 404.
	f, err := nfs.fs.Open(path)
	if err != nil {
		return nil, err
	}

	// It exists, see if it's a directory. That's the case we're interested in
	// for special handling.
	s, err := f.Stat()

	if err == nil && s.IsDir() {

		// If entry is a directory, check if index.html exists by attempting Open().

		// quick-and-dirty fix for windows platform: convert its "\" separator
		// to "/" that FileSystem demands for its paths. Should use URL join, tho.
		index := strings.Replace(filepath.Join(path, "index.html"), `\`, `/`, -1)

		if findex, err := nfs.fs.Open(index); err != nil {

			// Presumably a File Not Found error is the cause. Close the directory
			// and let caller think directory open is the source of Not Found err
			// we are returning.
			closeErr := f.Close()

			// shouldn't fail, but just in case, report the error whatever it is
			if closeErr != nil {
				return nil, closeErr
			}

			// Dir exists with no index file. Return err Not Found rather than
			// let dir listing be served
			return nil, err

		} else {
			// open test succeeded, so close it. No further use for it.
			findex.Close()
		}

		// directory has Openable index.html. Fall thru and let FileServer
		// do its normal thing: detect and return it to client instead of dir listing

	}

	return f, nil
}
