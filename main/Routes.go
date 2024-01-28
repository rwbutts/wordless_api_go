package main

import "net/http"

func MapRoutes(mux *http.ServeMux) {

	fs := http.FileServer(noListFileSystem{http.Dir("wwwroot")})

	mux.HandleFunc("/api/healthcheck", handleHealthCheck)
	mux.HandleFunc("/api/randomword", handleGetRandomWord)
	mux.HandleFunc("/api/getword/", handleGetWord)
	mux.HandleFunc("/api/checkword/", handleCheckWord)
	mux.HandleFunc("/api/querymatchcount", handleMatchCount)
	mux.Handle("/", fs)

}
