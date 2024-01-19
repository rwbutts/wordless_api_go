package main

import (
	"fmt"
	"net/http"
	"path/filepath"
	"strings"
)

func main() {

	config := GetConfig()

	fmt.Println("HTTP listening on " + config.ListenAddress)

	fs := http.FileServer(noListFileSystem{http.Dir("wwwroot")})

	http.HandleFunc("/api/healthcheck", handleHealthCheck)
	http.HandleFunc("/api/randomword", handleGetRandomWord)
	http.HandleFunc("/api/getword/", handleGetWord)
	http.HandleFunc("/api/checkword/", handleCheckWord)
	http.HandleFunc("/api/querymatchcount", handleMatchCount)
	http.Handle("/", fs)

	err := http.ListenAndServe(config.ListenAddress, nil)
	if err != nil {
		panic(err.Error())
	}

}

type noListFileSystem struct {
	fs http.FileSystem
}

/*
Open() functon for noListFileSystem struct that wraps the normal fs.Open functionality.
Inhibits the default FileServer "if it'a a directory without index.html,
show a directory listing" functionality.  With this, it will return a file not found
error, which results instead in a client 404 not found rather than serving dir listing.
*/
func (nfs noListFileSystem) Open(path string) (http.File, error) {
	/*
		If neither file nor directory, just return error so http returns 404.
	*/
	f, err := nfs.fs.Open(path)
	if err != nil {
		return nil, err
	}

	/*
		It exists, see if it's a directory. That's the case we're interested in
		for special handling.
	*/
	s, _ := f.Stat()

	/*
		If entry is a directory, check if index.html exists. If not, cause
		http to return a 494 error; if it does, just let http do what
		it normally does: return the index.html in that dir.
	*/
	if s.IsDir() {
		// quick-and-dirty fix for windows platorm: convert the "\" separator
		// to "/" that FileSystem requires for its paths. Should user URL join, tho.
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
