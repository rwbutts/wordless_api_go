package main

import (
	"net/http"
)

const LISTEN_ADDRESS = ":5090"

func main() {

	fs := http.FileServer(http.Dir("./wwwroot"))

	http.HandleFunc("/api/healthcheck", handleHealthCheck)
	http.HandleFunc("/api/randomword", handleGetRandomWord)
	http.HandleFunc("/api/getword/", handleGetWord)
	http.HandleFunc("/api/checkword/", handleCheckWord)
	http.HandleFunc("/api/querymatchcount", handleMatchCount)
	http.Handle("/", fs)

	err := http.ListenAndServe(LISTEN_ADDRESS, nil)
	if err != nil {
		panic(err.Error())
	}

}
