package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func main() {

	config := GetConfig()

	log.SetOutput(io.Discard)

	if config.LogFile != "" {
		logfile, err := os.OpenFile(config.LogFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
		if err != nil {
			fmt.Printf("Warning: error opening log file %v: %v", config.LogFile, err.Error())
		} else {
			defer logfile.Close()
			log.SetOutput(logfile)
		}
	}

	defer log.Println("Shutdown")

	fmt.Println("HTTP listening on " + config.ListenAddress)
	log.Println("HTTP listening on " + config.ListenAddress)

	mux := http.NewServeMux()

	fs := http.FileServer(noListFileSystem{http.Dir("wwwroot")})

	mux.Handle("/api/healthcheck", LogMiddleware(http.HandlerFunc(handleHealthCheck)))
	mux.Handle("/api/randomword", LogMiddleware(http.HandlerFunc(handleGetRandomWord)))
	mux.Handle("/api/getword/", LogMiddleware(http.HandlerFunc(handleGetWord)))
	mux.Handle("/api/checkword/", LogMiddleware(http.HandlerFunc(handleCheckWord)))
	mux.Handle("/api/querymatchcount", LogMiddleware(http.HandlerFunc(handleMatchCount)))
	mux.Handle("/", LogMiddleware(fs))

	err := http.ListenAndServe(config.ListenAddress, mux)
	if err != nil {
		log.Println("PANIC : " + err.Error())
		panic(err.Error())
	}
}

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
	s, _ := f.Stat()

	// If entry is a directory, check if index.html exists. If not, cause
	// http to return a 404 error; if it does, just let http do what
	// it normally does: return the index.html in that dir.
	if s.IsDir() {
		// quick-and-dirty fix for windows platform: convert its "\" separator
		// to "/" that FileSystem demands for its paths. Should use URL join, tho.
		index := strings.Replace(filepath.Join(path, "index.html"), `\`, `/`, -1)
		if _, err := nfs.fs.Open(index); err != nil {
			closeErr := f.Close()
			if closeErr != nil {
				return nil, closeErr
			}

			return nil, err
		}
	}

	return f, nil
}
