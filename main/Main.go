package main

import (
	"net/http"
)

const LISTEN_ADDRESS = ":5090"

func main() {

	http.HandleFunc("/api/healthcheck", handleHealthCheck)
	http.HandleFunc("/api/randomword", handleGetRandomWord)
	http.HandleFunc("/api/getword/", handleGetWord)
	http.HandleFunc("/api/checkword/", handleCheckWord)
	http.HandleFunc("/api/querymatchcount", handleMatchCount)

	err := http.ListenAndServe(LISTEN_ADDRESS, nil)
	if err != nil {
		panic(err.Error())
	}

}
