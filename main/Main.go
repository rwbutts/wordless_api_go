package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

const EXENAME = "wordless"

func main() {

	if len(os.Args) > 1 && (os.Args[1] == "--help" || os.Args[1] == "-h" || os.Args[1] == "/?") {
		fmt.Printf("Usage: %v [json config file]\n", EXENAME)
		fmt.Println("\nSample file:")
		fmt.Println("    {")
		fmt.Println(`		  "port" : "8080"`)
		fmt.Println(`		  "logfile" : "config.json"`)
		fmt.Println("    }")
		fmt.Println("\nConfig can also be set by the environment variables: PORT and LOGFILE.")
		fmt.Println("Note: Environment variable has higher priority than config file")
		fmt.Println("\nDefault port is 8080; default logfile = \"\" (no logging)")
		fmt.Println(`To log to stdout, set logfile =  "-"`)

		os.Exit(0)
	}

	config := GetConfig()

	if config.LogFile == "" {
		log.SetOutput(io.Discard)
	} else if config.LogFile == "-" {
		log.SetOutput(os.Stdout)
	} else {
		logfile, err := os.OpenFile(config.LogFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
		if err != nil {
			fmt.Printf("Warning: error opening log file %v: %v", config.LogFile, err.Error())
		} else {
			defer logfile.Close()
			log.SetOutput(logfile)
		}
	}

	defer log.Println("Shutdown")

	fmt.Println("HTTP listening on :" + config.Port)
	log.Println("HTTP listening on :" + config.Port)

	mux := http.NewServeMux()

	fs := http.FileServer(noListFileSystem{http.Dir("wwwroot")})

	mux.HandleFunc("/api/healthcheck", handleHealthCheck)
	mux.HandleFunc("/api/randomword", handleGetRandomWord)
	mux.HandleFunc("/api/getword/", handleGetWord)
	mux.HandleFunc("/api/checkword/", handleCheckWord)
	mux.HandleFunc("/api/querymatchcount", handleMatchCount)
	mux.Handle("/", fs)

	mw := LoggingMiddleware(log.Default())

	err := http.ListenAndServe(":"+config.Port, mw(mux))
	if err != nil {
		log.Println("PANIC ListenAndServe() exit: " + err.Error())
		panic(err.Error())
	}
}
